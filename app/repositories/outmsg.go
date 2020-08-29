package repositories

import (
	"message-queue/app/models"
	"message-queue/app/schema"
)

type OutRepository interface {
	Retrieve(id string) (*models.OutMessage, error)
	Get(query *schema.OutMsgQueryParam) (*models.OutMessage, error)
	List(query *schema.OutMsgQueryParam) (*[]models.OutMessage, error)
	Create(message *models.OutMessage) error
	Update(message *models.OutMessage) error
}
