package repositories

import (
	"github.com/quangdangfit/gosdk/utils/paging"

	"message-queue/app/models"
	"message-queue/app/schema"
)

type InRepository interface {
	Retrieve(id string) (*models.InMessage, error)
	Get(query *schema.InMsgQueryParam) (*models.InMessage, error)
	List(query *schema.InMsgQueryParam) (*[]models.InMessage, *paging.Paging, error)
	Create(message *models.InMessage) error
	Update(message *models.InMessage) error
	Upsert(message *models.InMessage) error
}
