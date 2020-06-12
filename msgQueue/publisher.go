package msgQueue

import (
	"encoding/json"
	"fmt"
	"gomq/config"
	"gomq/dbs"
	"gomq/models"
	"gomq/msgHandler"
	"transport/lib/utils/logger"

	"github.com/streadway/amqp"
)

type Publisher interface {
	MessageQueue
	Publish(message *models.OutMessage, reliable bool) error
	confirmAndHandle(message *models.OutMessage, confirms chan amqp.Confirmation) error
}

type publisher struct {
	messageQueue
}

func NewPublisher() Publisher {
	var pub publisher

	pub.config = &models.AMQPConfig{
		AMQPUrl:      config.Config.AMQP.URL,
		ExchangeName: config.Config.AMQP.ExchangeName,
		ExchangeType: config.Config.AMQP.ExchangeType,
	}
	pub.newConnection()
	return &pub
}

func (pub *publisher) Publish(message *models.OutMessage, reliable bool) (
	err error) {

	defer pub.closeConnection()

	// Reliable publisher confirms require confirm.select support from the connection.
	if reliable {
		if err := pub.channel.Confirm(false); err != nil {
			logger.Errorf(
				"Channel could not be put into confirm mode: %s", err)
			return err
		}
		confirms := pub.channel.NotifyPublish(make(chan amqp.Confirmation, 1))
		defer pub.confirmAndHandle(message, confirms)
	}

	payload, _ := json.Marshal(message.Payload)
	headers := amqp.Table{
		"origin_code":  message.OriginCode,
		"origin_model": message.OriginModel,
	}
	if err = pub.channel.Publish(
		pub.config.ExchangeName, // publish to an exchange
		message.RoutingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Headers:         headers,
			ContentType:     "application/json",
			ContentEncoding: "",
			Body:            payload,
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		message.Status = dbs.OutMessageStatusFailed
		return fmt.Errorf("Exchange Publish: %s", err)
	}

	return nil
}

func (pub *publisher) confirmAndHandle(message *models.OutMessage, confirms chan amqp.Confirmation) error {
	pub.confirmOne(message, confirms)

	outHandler := msgHandler.NewOutMessageHandler()
	_, err := outHandler.HandleMessage(message)
	return err
}
