package repositories

import (
	"gomq/app/models"
	"gomq/app/schema"
)

type InMessageRepository interface {
	GetInMessageByID(id string) (*models.InMessage, error)
	GetSingleInMessage(query *schema.InMessageQueryParam) (*models.InMessage, error)
	GetInMessages(query *schema.InMessageQueryParam, limit int) (*[]models.InMessage, error)
	AddInMessage(message *models.InMessage) error
	UpdateInMessage(message *models.InMessage) error
	UpsertInMessage(message *models.InMessage) error
}

type OutMessageRepository interface {
	GetOutMessageByID(id string) (*models.OutMessage, error)
	GetSingleOutMessage(query *schema.OutMessageQueryParam) (*models.OutMessage, error)
	GetOutMessages(query *schema.OutMessageQueryParam, limit int) (*[]models.OutMessage, error)
	AddOutMessage(message *models.OutMessage) error
	UpdateOutMessage(message *models.OutMessage) error
}

type RoutingRepository interface {
	GetRoutingKey(query map[string]interface{}) (*models.RoutingKey, error)
	GetPreviousRoutingKey(srcRouting models.RoutingKey) (*models.RoutingKey, error)
}
