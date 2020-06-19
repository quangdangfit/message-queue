package msgQueue

import (
	"gomq/config"
	"gomq/dbs"
	"gomq/models"

	"github.com/streadway/amqp"
	"gitlab.com/quangdangfit/gocommon/utils/logger"
)

const RecoverIntervalTime = 6 * 60
const TimeoutRetry = 3

type MessageQueue interface {
	newConnection() (*amqp.Connection, error)
	closeConnection() error
	ensureConnection() (err error)
	newChannel() (*amqp.Channel, error)
	closeChannel() error
	declareExchange() error
	declareQueue() error
	bindQueue(routingKey string) error
	confirmOne(message *models.OutMessage,
		confirms <-chan amqp.Confirmation) bool
	setup()
	ChanelIsClosed() bool
}

type messageQueue struct {
	config          *models.AMQPConfig
	connection      *amqp.Connection
	channel         *amqp.Channel
	errorChan       chan *amqp.Error
	isClosed        bool
	channelIsClosed bool
}

func (mq *messageQueue) newConnection() (*amqp.Connection, error) {
	conn, err := amqp.Dial(mq.config.AMQPUrl)
	if err != nil {
		logger.Error("Cannot create new connection to AMQP: ", err)
		return nil, err
	}
	mq.connection = conn
	mq.newChannel()

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
	channel, err := mq.connection.Channel()
	if err != nil {
		logger.Error("Failed to new channel: ", err)
		return nil, err
	}
	mq.channel = channel
	return channel, nil
}

func (mq *messageQueue) ensureConnection() (err error) {
	if mq.connection.IsClosed() {
		_, err = mq.newConnection()
	}
	return err
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
	if mq.ChanelIsClosed() {
		logger.Error("Channel is closed, cannot declare exchange")
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
	if mq.ChanelIsClosed() {
		logger.Error("Channel is closed, cannot declare exchange")
	}

	if _, err := mq.channel.QueueDeclare(
		mq.config.QueueName,
		false,
		false,
		true,
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

func (mq *messageQueue) confirmOne(message *models.OutMessage,
	confirms <-chan amqp.Confirmation) bool {

	confirmed := <-confirms
	if confirmed.Ack {
		logger.Info("Confirmed delivery with delivery tag: ",
			confirmed.DeliveryTag)

		message.Status = dbs.OutMessageStatusSent
		return true
	}

	logger.Info("Failed delivery of delivery tag: ",
		confirmed.DeliveryTag)

	message.Status = dbs.OutMessageStatusSentWait
	return false
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
