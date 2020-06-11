package msgQueue

import (
	"encoding/json"
	"gomq/config"
	"gomq/models"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Consumer interface {
	MessageQueue
	Consume() (<-chan models.InMessage, error)
	loopConsumeMsg(delivery <-chan amqp.Delivery, msgChan chan<- models.InMessage) (err error)
	channelConsume() (<-chan amqp.Delivery, error)
	convertAMQPMsgToMessage(amqpMsg amqp.Delivery) models.InMessage
}

type consumer struct {
	messageQueue
}

func NewConsumer() Consumer {
	var sub consumer

	sub.config = &models.AMQPConfig{
		AMQPUrl:   config.Config.AMQP.URL,
		QueueName: config.Config.AMQP.QueueName,
	}
	sub.NewConnection()
	sub.DeclareQueue()

	return &sub
}

func (cons *consumer) Consume() (<-chan models.InMessage, error) {
	defer cons.CloseConnection()

	delivery, err := cons.channelConsume()
	if err != nil {
		return nil, err
	}

	msgChan := make(chan models.InMessage)
	go cons.loopConsumeMsg(delivery, msgChan)
	return msgChan, nil
}

func (cons *consumer) loopConsumeMsg(delivery <-chan amqp.Delivery, msgChan chan<- models.InMessage) (err error) {
	for !cons.isClosed {
		select {
		case chanErr := <-cons.errorChan:
			logrus.Error("consume receive error: ", chanErr)
			logrus.Info("reconsume")
			for {
				delivery, err = cons.channelConsume()
				if err == nil {
					logrus.Info("reconsume success")
					break
				}

				logrus.Info("failed to reconsume, trying in 2s ...")
				time.Sleep(time.Second * 2)
			}

		case amqpMsg := <-delivery:
			msg := cons.convertAMQPMsgToMessage(amqpMsg)
			msgChan <- msg
		}
	}
	return nil
}

func (cons *consumer) channelConsume() (<-chan amqp.Delivery, error) {
	return cons.channel.Consume(
		cons.config.QueueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}

func (cons *consumer) convertAMQPMsgToMessage(amqpMsg amqp.Delivery) models.InMessage {
	var message models.InMessage
	json.Unmarshal(amqpMsg.Body, &message)
	return message
}
