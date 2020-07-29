package repositories

import (
	"go.uber.org/dig"

	"gomq/packages/models"
)

type InMessageRepository interface {
	GetSingleInMessage(query map[string]interface{}) (*models.InMessage, error)
	GetInMessages(query map[string]interface{}, limit int) (*[]models.InMessage, error)
	AddInMessage(message *models.InMessage) error
	UpdateInMessage(message *models.InMessage) error
	UpsertInMessage(message *models.InMessage) error
}

type OutMessageRepository interface {
	GetSingleOutMessage(query map[string]interface{}) (*models.OutMessage, error)
	GetOutMessages(query map[string]interface{}, limit int) (*[]models.OutMessage, error)
	AddOutMessage(message *models.OutMessage) error
	UpdateOutMessage(message *models.OutMessage) error
}

type RoutingRepository interface {
	GetRoutingKey(query map[string]interface{}) (*models.RoutingKey, error)
	GetPreviousRoutingKey(srcRouting models.RoutingKey) (*models.RoutingKey, error)
}

func Inject(container *dig.Container) {
	_ = container.Provide(NewInMessageRepository)
	_ = container.Provide(NewOutMessageRepository)
	_ = container.Provide(NewRoutingRepository)
}
