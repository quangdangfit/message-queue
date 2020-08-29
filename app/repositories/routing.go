package repositories

import (
	"github.com/quangdangfit/gosdk/utils/paging"

	"gomq/app/models"
	"gomq/app/schema"
)

type RoutingRepository interface {
	Retrieve(query map[string]interface{}) (*models.RoutingKey, error)
	GetPrevious(srcRouting models.RoutingKey) (*models.RoutingKey, error)
	List(query *schema.RoutingKeyQueryParam) (*[]models.RoutingKey, *paging.Paging, error)
}
