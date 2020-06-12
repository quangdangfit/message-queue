package msgQueue

import (
	"errors"
	"fmt"
	"github.com/manucorporat/try"
	"gomq/config"
	"gomq/models"
	"gomq/msgHandler"
	"runtime"
	"sync/atomic"
	"time"
	"transport/lib/utils/logger"

	"github.com/streadway/amqp"
)

const RecoverIntervalTime = 6 * 60
const TimeoutRetry = 3

type Consumer interface {
	MessageQueue
	RunConsumer(handler func([]byte) bool)
	reconnect(retryTime int) (<-chan amqp.Delivery, error)
	announceQueue() (<-chan amqp.Delivery, error)
	handle(deliveries <-chan amqp.Delivery, fn func([]byte) bool, threads int)
	consume(deliveries <-chan amqp.Delivery)
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
		AMQPUrl:   config.Config.AMQP.URL,
		QueueName: config.Config.AMQP.QueueName,
	}
	sub.newConnection()
	sub.declareQueue()

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
	if err := cons.newConnection(); err != nil {
		fmt.Sprintf("[%s]connect error", "consumerTag")
	}

	deliveries, _ := cons.announceQueue()
	fmt.Sprintf("[%s]Error when calling announceQueue()", "consumerTag")
	//cons.handle(deliveries, handler, maxParallelism(), queueName, routingKey)
	cons.handle(deliveries, handler, 3)
}

func (cons *consumer) reconnect(retryTime int) (<-chan amqp.Delivery, error) {
	cons.closeConnection()
	time.Sleep(time.Duration(TimeoutRetry) * time.Second)
	logger.Info("Try reConnect with times:", retryTime)

	if err := cons.newConnection(); err != nil {
		return nil, err
	}

	deliveries, err := cons.announceQueue()
	if err != nil {
		return deliveries, errors.New("Couldn't connect")
	}
	return deliveries, nil
}

// announceQueue sets the queue that will be listened to for this connection
func (cons *consumer) announceQueue() (<-chan amqp.Delivery, error) {
	err := cons.channel.Qos(50, 0, false)
	if err != nil {
		return nil, fmt.Errorf("Error setting qos: %s", err)
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
		return nil, fmt.Errorf("Queue Consume: %s", err)
	}
	return deliveries, nil
}

func (cons *consumer) handle(deliveries <-chan amqp.Delivery,
	fn func([]byte) bool, threads int) {
	var err error
	for {
		logger.Info("Enter for busy loop with thread:", threads)
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
	logger.Info("Enter go with thread with deliveries", deliveries)
	for msg := range deliveries {
		logger.Info("Enter deliver")
		ret := false
		try.This(func() {
			receiver := msgHandler.NewReceiver()
			message, err := receiver.HandleMessage(msg.Body, msg.RoutingKey)
			if err != nil {
				logger.Errorf("Failed to handle message, routing_key %s, "+
					"model %s, code %s", message.RoutingKey,
					message.OriginModel, message.OriginCode)
			}
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
