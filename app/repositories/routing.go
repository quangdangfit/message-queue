package repositories

import (
	"github.com/quangdangfit/gosdk/utils/paging"

	"message-queue/app/models"
	"message-queue/app/schema"
)

type RoutingRepository interface {
	Get(query map[string]interface{}) (*models.RoutingKey, error)
	GetPrevious(srcRouting models.RoutingKey) (*models.RoutingKey, error)
	List(query *schema.RoutingQueryParam) (*[]models.RoutingKey, *paging.Paging, error)
	Create(body *schema.RoutingCreateParam) (*models.RoutingKey, error)
	Retrieve(id string) (*models.RoutingKey, error)
}
