package queue

import (
	"encoding/json"

	"github.com/quangdangfit/gosdk/utils/logger"
	"github.com/streadway/amqp"

	"gomq/app/models"
	"gomq/config"
	"gomq/pkg/utils"
)

type Publisher interface {
	Publish(message *models.OutMessage, reliable bool) error
	confirmOne(message *models.OutMessage, confirms <-chan amqp.Confirmation) bool
}

type publisher struct {
	messageQueue
}

func NewPublisher() Publisher {
	pub := publisher{}
	pub.config = &AMQPConfig{
		AMQPUrl:      config.Config.AMQP.URL,
		ExchangeName: config.Config.AMQP.ExchangeName,
		ExchangeType: config.Config.AMQP.ExchangeType,
		QueueName:    config.Config.AMQP.QueueName,
	}
	_, err := pub.newConnection()
	if err != nil {
		logger.Error("Publisher create new connection failed!")
	}

	err = pub.declareExchange()
	if err != nil {
		logger.Error("Publisher declare exchange failed!")
	}

	return &pub
}

func (pub *publisher) Publish(message *models.OutMessage, reliable bool) (
	err error) {

	// New channel and close after publish
	pub.ensureConnection()
	channel, _ := pub.connection.Channel()
	defer channel.Close()

	// Reliable publisher confirms require confirm.select support from the connection.
	var confirms chan amqp.Confirmation
	if reliable {
		if err := channel.Confirm(false); err != nil {
			logger.Errorf("Channel could not be put into confirm mode: %s", err)
			return err
		}
		confirms = channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	}

	payload, _ := json.Marshal(message.Payload)
	headers := amqp.Table{
		"origin_code":  message.OriginCode,
		"origin_model": message.OriginModel,
		"api_key":      message.APIKey,
	}
	if err = channel.Publish(
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
		message.Status = models.OutMessageStatusFailed
		message.Logs = append(message.Logs, utils.ParseLogs(err))
		logger.Error("Failed to publish message ", err)
		return err
	}

	if confirms != nil {
		defer pub.confirmOne(message, confirms)
	}

	return nil
}

func (pub *publisher) confirmOne(message *models.OutMessage, confirms <-chan amqp.Confirmation) bool {

	confirmed := <-confirms
	if confirmed.Ack {
		logger.Info("Confirmed delivery with delivery tag: ",
			confirmed.DeliveryTag)

		message.Status = models.OutMessageStatusSent
		return true
	}

	logger.Info("Failed delivery of delivery tag: ",
		confirmed.DeliveryTag)

	message.Status = models.OutMessageStatusSentWait
	return false
}
