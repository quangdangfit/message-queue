package queue

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"gitlab.com/quangdangfit/gocommon/utils/logger"

	"gomq/config"
)

const (
	RecoverIntervalTime = 6 * 60
	TimeoutRetry        = 3
	WaitTimeReconnect   = 5
)

type MessageQueue interface {
	newConnection() (*amqp.Connection, error)
	closeConnection() error
	ensureConnection() (err error)
	newChannel() (*amqp.Channel, error)
	closeChannel() error
	declareExchange() error
	declareQueue() error
	bindQueue(routingKey string) error
	ChanelIsClosed() bool
}

type messageQueue struct {
	config          *AMQPConfig
	connection      *amqp.Connection
	channel         *amqp.Channel
	errorChan       chan *amqp.Error
	isClosed        bool
	channelIsClosed bool
}

func (mq *messageQueue) newConnection() (*amqp.Connection, error) {
	conn, err := amqp.Dial(mq.config.AMQPUrl)
	for err != nil {
		logger.Error("Failed to create new connection to AMQP: ", err)

		logger.Infof("Sleep %d seconds to reconnect", WaitTimeReconnect)
		time.Sleep(WaitTimeReconnect * time.Second)
		conn, err = amqp.Dial(mq.config.AMQPUrl)
	}
	mq.connection = conn

	return conn, nil
}

func (mq *messageQueue) closeConnection() error {
	if mq.isClosed {
		return nil
	}
	mq.closeChannel()

	if mq.connection != nil {
		if err := mq.connection.Close(); err != nil {
			return err
		}
		mq.connection = nil
	}

	mq.isClosed = true
	return nil
}

func (mq *messageQueue) newChannel() (*amqp.Channel, error) {
	mq.ensureConnection()

	if mq.connection == nil || mq.connection.IsClosed() {
		logger.Error("Connection is not open, cannot create new channel")
		return nil, fmt.Errorf("Connection is not open")
	}

	channel, err := mq.connection.Channel()
	if err != nil {
		logger.Error("Failed to new channel: ", err)
		return nil, err
	}
	mq.channel = channel
	mq.channelIsClosed = false
	logger.Info("New channel successfully")
	return channel, nil
}

func (mq *messageQueue) ensureConnection() (err error) {
	if mq.connection == nil || mq.connection.IsClosed() {
		_, err = mq.newConnection()
		if err != nil {
			return err
		}
	}
	return nil
}

func (mq *messageQueue) closeChannel() error {
	if mq.isClosed || mq.channelIsClosed {
		return nil
	}
	logger.Info("Close channel")
	if mq.channel != nil {
		_ = mq.channel.Close()
		mq.channel = nil
		mq.channelIsClosed = true
	}

	return nil
}

func (mq *messageQueue) declareExchange() error {
	mq.newChannel()
	defer mq.closeChannel()

	if mq.ChanelIsClosed() {
		logger.Error("Channel is not open, cannot declare exchange")
	}

	if err := mq.channel.ExchangeDeclare(
		mq.config.ExchangeName, // name
		mq.config.ExchangeType, // type
		true,                   // durable
		false,                  // auto-deleted
		false,                  // internal
		false,                  // noWait
		nil,                    // arguments
	); err != nil {
		logger.Error("Failed to declare exchange: ", err)
		return err
	}

	logger.Info("Declared exchange: ", mq.config.ExchangeName)
	return nil
}

func (mq *messageQueue) declareQueue() error {
	mq.newChannel()
	defer mq.closeChannel()

	if mq.ChanelIsClosed() {
		logger.Error("Channel is not open, cannot declare exchange")
	}

	if _, err := mq.channel.QueueDeclare(
		mq.config.QueueName,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		logger.Error("Failed to declare queue: ", err)
		return err
	}

	logger.Info("Declared queue: ", mq.config.QueueName)
	return nil
}

func (mq *messageQueue) bindQueue(routingKey string) error {
	if err := mq.channel.QueueBind(
		mq.config.QueueName,             // name
		routingKey,                      // key
		config.Config.AMQP.ExchangeName, // exchange
		false,                           //noWait
		nil,                             // args
	); err != nil {
		logger.Error("Failed to bind queue: ", err)
		return err
	}
	return nil
}

func (mq *messageQueue) setup() {
	mq.declareExchange()
	mq.declareQueue()
}

func (mq *messageQueue) ChanelIsClosed() bool {
	if mq.channel == nil || mq.channelIsClosed {
		return true
	}
	return false
}
