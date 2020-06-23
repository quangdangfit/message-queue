package msgQueue

import (
	"encoding/json"
	"errors"
	"github.com/manucorporat/try"
	"gomq/config"
	"gomq/models"
	"gomq/msgHandler"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/streadway/amqp"
	"gitlab.com/quangdangfit/gocommon/utils/logger"
)

const (
	ConsumerThreads = 10
)

type Consumer interface {
	MessageQueue
	RunConsumer(handler func([]byte) bool)
	getHeaderKey(msg amqp.Delivery, key string) string
	parseMessageFromDelivery(msg amqp.Delivery) (*models.InMessage, error)
	reconnect(retryTime int) (<-chan amqp.Delivery, error)
	subscribe() (<-chan amqp.Delivery, error)
	startConsuming(deliveries <-chan amqp.Delivery, fn func([]byte) bool, threads int)
	consume(deliveries <-chan amqp.Delivery)
	handle(msg amqp.Delivery)
}

type consumer struct {
	messageQueue

	done            chan error
	consumerTag     string // Name that consumer identifies itself to the server
	lastRecoverTime int64
	//track service current status
	currentStatus atomic.Value
}

func NewConsumer() Consumer {
	var sub = consumer{
		done:            make(chan error),
		lastRecoverTime: time.Now().Unix(),
	}

	sub.config = &models.AMQPConfig{
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

	return &sub
}

func maxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU

}

func (cons *consumer) RunConsumer(handler func([]byte) bool) {
	cons.newChannel()
	deliveries, _ := cons.subscribe()

	threads := config.Config.AMQP.ConsumerThreads
	if threads <= 0 {
		threads = ConsumerThreads
	}
	cons.startConsuming(deliveries, handler, threads)
}

func (cons *consumer) getHeaderKey(msg amqp.Delivery, key string) string {
	if msg.Headers[key] != nil {
		return msg.Headers[key].(string)
	}
	return ""
}

func (cons *consumer) parseMessageFromDelivery(msg amqp.Delivery) (
	*models.InMessage, error) {

	var payload interface{}
	json.Unmarshal(msg.Body, &payload)
	message := models.InMessage{
		Payload:     payload,
		OriginCode:  msg.Headers["origin_code"].(string),
		OriginModel: msg.Headers["origin_model"].(string),
		APIKey:      cons.getHeaderKey(msg, "api_key"),
	}
	return &message, nil
}

func (cons *consumer) convertMessage(bytesMsg []byte, payload interface{}) error {
	json.Unmarshal(bytesMsg, &payload)
	return nil
}

func (cons *consumer) reconnect(retryTime int) (<-chan amqp.Delivery, error) {
	cons.closeConnection()
	time.Sleep(time.Duration(TimeoutRetry) * time.Second)
	logger.Info("Try reConnect with times:", retryTime)

	cons.ensureConnection()

	deliveries, err := cons.subscribe()
	if err != nil {
		return deliveries, errors.New("Couldn't connect")
	}
	return deliveries, nil
}

// subscribe sets the queue that will be listened to for this connection
func (cons *consumer) subscribe() (<-chan amqp.Delivery, error) {
	err := cons.channel.Qos(50, 0, false)
	if err != nil {
		logger.Error("Error setting qos: ", err)
		return nil, err
	}

	logger.Info("Queue bound to Exchange, starting Consume consumer tag:",
		cons.consumerTag)

	deliveries, err := cons.channel.Consume(
		cons.config.QueueName, // name
		cons.consumerTag,      // consumerTag,
		false,                 // noAck
		false,                 // exclusive
		false,                 // noLocal
		false,                 // noWait
		nil,                   // arguments
	)
	if err != nil {
		logger.Error("Failed to consume queue: ", err)
		return nil, err
	}
	return deliveries, nil
}

func (cons *consumer) startConsuming(deliveries <-chan amqp.Delivery,
	fn func([]byte) bool, threads int) {
	var err error
	for {
		logger.Info("Enter for busy loop with threads: ", threads)
		for i := 0; i < threads; i++ {
			go cons.consume(deliveries)
		}

		// Go into reconnect loop when
		// cons.done is passed non nil values
		if <-cons.done != nil {
			cons.currentStatus.Store(false)
			retryTime := 1
			for {
				deliveries, err = cons.reconnect(retryTime)
				if err != nil {
					logger.Error("Reconnecting Error")
					retryTime += 1
				} else {
					break
				}
			}
		}
		logger.Info("Reconnected!!!")
	}
}

func (cons *consumer) consume(deliveries <-chan amqp.Delivery) {
	logger.Info("Enter go with thread with deliveries ", deliveries)
	for msg := range deliveries {
		logger.Info("Enter deliver")
		ret := false
		try.This(func() {
			go cons.handle(msg)
		}).Finally(func() {
			if ret == true {
				msg.Ack(false)
				currentTime := time.Now().Unix()
				if currentTime-cons.lastRecoverTime > RecoverIntervalTime &&
					!cons.currentStatus.Load().(bool) {

					logger.Info("Try to Recover Unack Messages!")
					cons.currentStatus.Store(true)
					cons.lastRecoverTime = currentTime
					cons.channel.Recover(true)
				}

			} else {
				// this really a litter dangerous. if the worker is panic very
				//quickly, it will ddos our sentry server......plz,
				//add [retry-ttl] in header.
				//msg.Nack(false, true)
				msg.Reject(false)
				//cons.currentStatus.Store(true)
			}
		}).Catch(func(e try.E) {
			logger.Error(e)
		})
	}
}

func (cons *consumer) handle(msg amqp.Delivery) {
	message, _ := cons.parseMessageFromDelivery(msg)
	receiver := msgHandler.NewInMessageHandler(true)
	_, err := receiver.HandleMessage(message, msg.RoutingKey)
	if err != nil {
		logger.Errorf("Failed to handle message, routing_key %s, "+
			"model %s, code %s, error: %s", message.RoutingKey.Name,
			message.OriginModel, message.OriginCode, err)
	}
}
