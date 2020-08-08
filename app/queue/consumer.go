package queue

import (
	"encoding/json"
	"errors"
	"sync/atomic"
	"time"

	"github.com/jinzhu/copier"
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
	Consume() chan *models.InMessage
	parseMessageFromDelivery(msg amqp.Delivery) (*models.InMessage, error)
	reconnect(retryTime int) (<-chan amqp.Delivery, error)
	subscribe() (<-chan amqp.Delivery, error)
	startConsuming(deliveries <-chan amqp.Delivery)
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

func (c *consumer) Consume() chan *models.InMessage {
	c.ensureConnection()
	c.newChannel()
	deliveries, _ := c.subscribe()
	go c.startConsuming(deliveries)

	// retry consume failed messages
	go c.startConsuming(c.failedChan)
	return c.msgChan
}

func (c *consumer) parseMessageFromDelivery(msg amqp.Delivery) (*models.InMessage, error) {
	var payload interface{}
	json.Unmarshal(msg.Body, &payload)
	var headers models.Headers
	data, err := json.Marshal(msg.Headers)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &headers)

	message := models.InMessage{
		Payload: payload,
	}
	copier.Copy(&message, &headers)
	message.RoutingKey.Name = msg.RoutingKey
	return &message, nil
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

func (c *consumer) startConsuming(deliveries <-chan amqp.Delivery) {
	logger.Info("Enter with deliveries ", deliveries)
	for msg := range deliveries {
		logger.Info("Enter deliver message: ", msg.RoutingKey)
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
	message, err := c.parseMessageFromDelivery(msg)
	if err != nil {
		c.failedChan <- msg
		logger.Error("Failed to parse message: ", err)
	}

	c.msgChan <- message
}
