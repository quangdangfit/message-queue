package services

import (
	"context"

	"github.com/quangdangfit/gosdk/utils/paging"

	"message-queue/app/models"
	"message-queue/app/schema"
)

type InService interface {
	Consume()
	List(ctx context.Context, query *schema.InMsgQueryParam) (*[]models.InMessage, *paging.Paging, error)
	CronRetry() error
	CronRetryPrevious() error
}
