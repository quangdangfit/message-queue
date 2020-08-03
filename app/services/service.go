package services

import (
	"net/http"

	"go.uber.org/dig"

	"gomq/app/models"
	"gomq/app/schema"
)

type InMessageService interface {
	HandleMessage(message *models.InMessage, routingKey string) error
	storeMessage(message *models.InMessage) (err error)
	callAPI(message *models.InMessage) (*http.Response, error)
}

type OutMessageService interface {
	HandleMessage(message *models.OutMessage) (err error)
	GetOutMessages(query *schema.OutMessageQueryParam, limit int) (*[]models.OutMessage, error)
}

func Inject(container *dig.Container) error {
	_ = container.Provide(NewInMessageService)
	_ = container.Provide(NewOutMessageService)

	return nil
}
