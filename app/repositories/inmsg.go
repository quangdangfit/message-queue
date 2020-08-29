package repositories

import (
	"message-queue/app/models"
	"message-queue/app/schema"
)

type InRepository interface {
	GetByID(id string) (*models.InMessage, error)
	Retrieve(query *schema.InMessageQueryParam) (*models.InMessage, error)
	List(query *schema.InMessageQueryParam, limit int) (*[]models.InMessage, error)
	Create(message *models.InMessage) error
	Update(message *models.InMessage) error
	Upsert(message *models.InMessage) error
}
