package repositories

import (
	"github.com/quangdangfit/gosdk/utils/paging"

	"gomq/app/models"
	"gomq/app/schema"
)

type RoutingRepository interface {
	GetRoutingKey(query map[string]interface{}) (*models.RoutingKey, error)
	GetPreviousRoutingKey(srcRouting models.RoutingKey) (*models.RoutingKey, error)
	GetRoutingKeys(query *schema.RoutingKeyQueryParam) (*[]models.RoutingKey, *paging.Paging, error)
}
