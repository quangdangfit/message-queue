package services

import (
	"context"

	"github.com/quangdangfit/gosdk/utils/paging"

	"message-queue/app/models"
	"message-queue/app/schema"
)

type RoutingService interface {
	Retrieve(ctx context.Context, id string) (*models.RoutingKey, error)
	List(ctx context.Context, query *schema.RoutingQueryParam) (*[]models.RoutingKey, *paging.Paging, error)
	Create(ctx context.Context, body *schema.RoutingCreateParam) (*models.RoutingKey, error)
}
