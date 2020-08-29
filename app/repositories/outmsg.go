package repositories

import (
	"message-queue/app/models"
	"message-queue/app/schema"
)

type OutRepository interface {
	GetByID(id string) (*models.OutMessage, error)
	Retrieve(query *schema.OutMsgQueryParam) (*models.OutMessage, error)
	List(query *schema.OutMsgQueryParam, limit int) (*[]models.OutMessage, error)
	Create(message *models.OutMessage) error
	Update(message *models.OutMessage) error
}
