package repositories

import (
	"gomq/app/models"
	"gomq/app/schema"
)

type OutRepository interface {
	GetByID(id string) (*models.OutMessage, error)
	Retrieve(query *schema.OutMessageQueryParam) (*models.OutMessage, error)
	List(query *schema.OutMessageQueryParam, limit int) (*[]models.OutMessage, error)
	Create(message *models.OutMessage) error
	Update(message *models.OutMessage) error
}
