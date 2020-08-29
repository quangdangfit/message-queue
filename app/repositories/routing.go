package repositories

import (
	"github.com/quangdangfit/gosdk/utils/paging"

	"message-queue/app/models"
	"message-queue/app/schema"
)

type RoutingRepository interface {
	Retrieve(id string) (*models.RoutingKey, error)
	Get(query *schema.RoutingQueryParam) (*models.RoutingKey, error)
	List(query *schema.RoutingQueryParam) (*[]models.RoutingKey, *paging.Paging, error)
	Create(body *schema.RoutingCreateParam) (*models.RoutingKey, error)
	Update(id string, body *schema.RoutingUpdateParam) (*models.RoutingKey, error)
}
