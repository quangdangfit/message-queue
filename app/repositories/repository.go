package repositories

import (
	"gomq/app/models"
	"gomq/app/schema"
)

type InRepository interface {
	GetInMessageByID(id string) (*models.InMessage, error)
	GetSingleInMessage(query *schema.InMessageQueryParam) (*models.InMessage, error)
	GetInMessages(query *schema.InMessageQueryParam, limit int) (*[]models.InMessage, error)
	AddInMessage(message *models.InMessage) error
	UpdateInMessage(message *models.InMessage) error
	UpsertInMessage(message *models.InMessage) error
}
