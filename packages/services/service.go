package services

import (
	"net/http"

	"go.uber.org/dig"

	"gomq/packages/models"
)

type InMessageService interface {
	HandleMessage(message *models.InMessage, routingKey string) error
	storeMessage(message *models.InMessage) (err error)
	callAPI(message *models.InMessage) (*http.Response, error)
}

type OutMessageService interface {
	HandleMessage(message *models.OutMessage) (err error)
}

func Inject(container *dig.Container) {
	_ = container.Provide(NewInMessageService)
	_ = container.Provide(NewOutMessageService)
}
