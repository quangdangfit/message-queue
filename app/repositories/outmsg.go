package repositories

import (
	"github.com/quangdangfit/gosdk/utils/paging"

	"message-queue/app/models"
	"message-queue/app/schema"
)

type OutRepository interface {
	Retrieve(id string) (*models.OutMessage, error)
	Get(query *schema.OutMsgQueryParam) (*models.OutMessage, error)
	List(query *schema.OutMsgQueryParam) (*[]models.OutMessage, *paging.Paging, error)
	Create(message *models.OutMessage) error
	Update(id string, body *schema.OutMsgUpdateParam) (*models.OutMessage, error)
}
