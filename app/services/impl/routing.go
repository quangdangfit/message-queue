package impl

import (
	"context"

	"github.com/quangdangfit/gosdk/utils/logger"
	"github.com/quangdangfit/gosdk/utils/paging"

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

func (r *routing) Retrieve(ctx context.Context, id string) (*models.RoutingKey, error) {
	rs, err := r.repo.Retrieve(id)
	if err != nil {
		logger.Errorf("Cannot get routing key %s, error: %s", id, err)
		return nil, err
	}

	return rs, nil
}

func (r *routing) List(ctx context.Context, query *schema.RoutingQueryParam) (*[]models.RoutingKey, *paging.Paging, error) {
	rs, pageInfo, err := r.repo.List(query)
	if err != nil {
		logger.Errorf("Cannot get list routing key, error: %s", err)
		return nil, nil, err
	}

	return rs, pageInfo, nil
}

func (r *routing) Create(ctx context.Context, body *schema.RoutingCreateParam) (*models.RoutingKey, error) {
	rs, err := r.repo.Create(body)
	if err != nil {
		logger.Error("Cannot create routing key, error: ", err)
		return nil, err
	}

	return rs, nil
}
