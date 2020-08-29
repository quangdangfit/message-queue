package impl

import (
	"context"

	"github.com/quangdangfit/gosdk/utils/logger"

	"message-queue/app/models"
	"message-queue/app/repositories"
	"message-queue/app/schema"
	"message-queue/app/services"
)

type routing struct {
	repo repositories.RoutingRepository
}

func NewRoutingService(repo repositories.RoutingRepository) services.RoutingService {
	return &routing{
		repo: repo,
	}
}

func (r *routing) Create(ctx context.Context, body *schema.RoutingCreateParam) (*models.RoutingKey, error) {
	rs, err := r.repo.Create(body)
	if err != nil {
		logger.Error("Cannot create routing key: ", err)
		return nil, err
	}

	return rs, nil
}
