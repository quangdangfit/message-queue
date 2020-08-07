package queue

import (
	"encoding/json"
	"errors"
	"sync/atomic"
	"time"

	"github.com/manucorporat/try"
	"github.com/quangdangfit/gosdk/utils/logger"
	"github.com/streadway/amqp"

	"gomq/app/models"
	"gomq/config"
)

const (
	DefaultConsumerThreads = 10
)

type Consumer interface {
	Consume(handler func([]byte) bool) chan *models.InMessage
	GetMessageChannel() chan *models.InMessage
	getHeaderKey(msg amqp.Delivery, key string) string
	parseMessageFromDelivery(msg amqp.Delivery) (*models.InMessage, error)
	reconnect(retryTime int) (<-chan amqp.Delivery, error)
	subscribe() (<-chan amqp.Delivery, error)
	startConsuming(deliveries <-chan amqp.Delivery, fn func([]byte) bool)
	pushToMsgChan(msg amqp.Delivery)
}

type consumer struct {
	messageQueue

	done            chan error
	consumerTag     string // Name that consumer identifies itself to the server
	lastRecoverTime int64
	//track service current status
	currentStatus atomic.Value

	threads int
	msgChan chan *models.InMessage
}

func NewConsumer() Consumer {
	var sub = consumer{
		done:            make(chan error),
		lastRecoverTime: time.Now().Unix(),
	}

	sub.config = &AMQPConfig{
		AMQPUrl:      config.Config.AMQP.URL,
		QueueName:    config.Config.AMQP.QueueName,
		ExchangeName: config.Config.AMQP.ExchangeName,
		ExchangeType: config.Config.AMQP.ExchangeType,
	}
	_, err := sub.newConnection()
	if err != nil {
		logger.Error("Consumer create new connection failed!")
	}

	err = sub.declareQueue()
	if err != nil {
		logger.Error("Consumer declare queue failed!")
	}

	threads := config.Config.AMQP.ConsumerThreads
	if threads <= 0 {
		threads = DefaultConsumerThreads
	}

	sub.threads = threads
	sub.msgChan = make(chan *models.InMessage, sub.threads)

	return &sub
}

func (c *consumer) Consume(handler func([]byte) bool) chan *models.InMessage {
	c.ensureConnection()
	c.newChannel()
	deliveries, _ := c.subscribe()
	go c.startConsuming(deliveries, handler)
	return c.msgChan
}

func (c *consumer) GetMessageChannel() chan *models.InMessage {
	return c.msgChan
}

func (c *consumer) getHeaderKey(msg amqp.Delivery, key string) string {
	if msg.Headers[key] != nil {
		return msg.Headers[key].(string)
	}
	return ""
}

func (c *consumer) parseMessageFromDelivery(msg amqp.Delivery) (
	*models.InMessage, error) {

	var payload interface{}
	json.Unmarshal(msg.Body, &payload)
	message := models.InMessage{
		Payload:     payload,
		OriginCode:  msg.Headers["origin_code"].(string),
		OriginModel: msg.Headers["origin_model"].(string),
		APIKey:      c.getHeaderKey(msg, "api_key"),
	}
	message.RoutingKey.Name = msg.RoutingKey
	return &message, nil
}

func (c *consumer) convertMessage(bytesMsg []byte, payload interface{}) error {
	json.Unmarshal(bytesMsg, &payload)
	return nil
}

func (c *consumer) reconnect(retryTime int) (<-chan amqp.Delivery, error) {
	c.closeConnection()
	time.Sleep(time.Duration(TimeoutRetry) * time.Second)
	logger.Info("Try reConnect with times:", retryTime)

	c.ensureConnection()

	deliveries, err := c.subscribe()
	if err != nil {
		return deliveries, errors.New("cannot connect")
	}
	return deliveries, nil
}

// subscribe sets the queue that will be listened to for this connection
func (c *consumer) subscribe() (<-chan amqp.Delivery, error) {
	err := c.channel.Qos(50, 0, false)
	if err != nil {
		logger.Error("Error setting qos: ", err)
		return nil, err
	}

	logger.Info("Queue bound to Exchange, starting Consume consumer tag:",
		c.consumerTag)

	deliveries, err := c.channel.Consume(
		c.config.QueueName, // name
		c.consumerTag,      // consumerTag,
		false,              // noAck
		false,              // exclusive
		false,              // noLocal
		false,              // noWait
		nil,                // arguments
	)
	if err != nil {
		logger.Error("Failed to consume queue: ", err)
		return nil, err
	}
	return deliveries, nil
}

func (c *consumer) startConsuming(deliveries <-chan amqp.Delivery, fn func([]byte) bool) {
	logger.Info("Enter with deliveries ", deliveries)
	for msg := range deliveries {
		logger.Info("Enter deliver")
		ret := false
		try.This(func() {
			c.pushToMsgChan(msg)
		}).Finally(func() {
			if ret == true {
				msg.Ack(false)
				currentTime := time.Now().Unix()
				if currentTime-c.lastRecoverTime > RecoverIntervalTime &&
					!c.currentStatus.Load().(bool) {

					logger.Info("Try to Recover Unack Messages!")
					c.currentStatus.Store(true)
					c.lastRecoverTime = currentTime
					c.channel.Recover(true)
				}

			} else {
				// this really a litter dangerous. if the worker is panic very
				//quickly, it will ddos our sentry server......plz,
				//add [retry-ttl] in header.
				//msg.Nack(false, true)
				msg.Reject(false)
				//c.currentStatus.Store(true)
			}
		}).Catch(func(e try.E) {
			logger.Error(e)
		})
	}
}

func (c *consumer) pushToMsgChan(msg amqp.Delivery) {
	message, _ := c.parseMessageFromDelivery(msg)
	c.msgChan <- message
}
