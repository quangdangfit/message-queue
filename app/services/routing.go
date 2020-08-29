package services

import (
	"context"

	"message-queue/app/models"
	"message-queue/app/schema"
)

type RoutingService interface {
	Create(ctx context.Context, body *schema.RoutingCreateParam) (*models.RoutingKey, error)
}
