package msgQueue

import (
	"gomq/config"
	"gomq/dbs"
	"gomq/models"
	"transport/lib/utils/logger"

	"github.com/streadway/amqp"
)

type MessageQueue interface {
	NewConnection() error
	CloseConnection() error
	NewChannel() error
	CloseChannel() error
	DeclareExchange() error
	DeclareQueue() error
	BindQueue(routingKey string) error
	ConfirmOne(message *models.OutMessage, confirms <-chan amqp.Confirmation) bool
}

type messageQueue struct {
	config     *models.AMQPConfig
	connection *amqp.Connection
	channel    *amqp.Channel
	errorChan  chan *amqp.Error
	isClosed   bool
}

func (mq *messageQueue) NewConnection() error {
	conn, err := amqp.Dial(mq.config.AMQPUrl)
	if err != nil {
		return err
	}
	mq.connection = conn
	mq.NewChannel()
	return nil
}

func (mq *messageQueue) CloseConnection() error {
	if mq.isClosed {
		return nil
	}
	mq.CloseChannel()

	if mq.connection != nil {
		if err := mq.connection.Close(); err != nil {
			return err
		}
		mq.connection = nil
	}
	return nil
}

func (mq *messageQueue) NewChannel() error {
	channel, err := mq.connection.Channel()
	if err != nil {
		logger.Error("Failed to new connection: ", err)
		return err
	}
	mq.channel = channel
	return nil
}

func (mq *messageQueue) CloseChannel() error {
	if mq.isClosed {
		return nil
	}
	logger.Info("Close channel")
	mq.isClosed = true
	if mq.channel != nil {
		_ = mq.channel.Close()
		mq.channel = nil
	}

	return nil
}

func (mq *messageQueue) DeclareExchange() error {
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
	return nil
}

func (mq *messageQueue) DeclareQueue() error {
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
	return nil
}

func (mq *messageQueue) BindQueue(routingKey string) error {
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

func (mq *messageQueue) ConfirmOne(message *models.OutMessage, confirms <-chan amqp.Confirmation) bool {
	confirmed := <-confirms
	if confirmed.Ack {
		logger.Info("Confirmed delivery with delivery tag: ", confirmed.DeliveryTag)
		message.Status = dbs.OutMessageStatusSent
		return true
	}

	logger.Info("Failed delivery of delivery tag: ", confirmed.DeliveryTag)
	message.Status = dbs.OutMessageStatusSentWait
	return false
}
