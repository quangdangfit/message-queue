package services

import (
	"context"

	"message-queue/app/models"
	"message-queue/app/schema"
)

type RoutingService interface {
	Retrieve(ctx context.Context, id string) (*models.RoutingKey, error)
	Create(ctx context.Context, body *schema.RoutingCreateParam) (*models.RoutingKey, error)
}
