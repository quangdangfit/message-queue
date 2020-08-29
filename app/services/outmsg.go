package services

import (
	"context"

	"github.com/quangdangfit/gosdk/utils/paging"

	"message-queue/app/models"
	"message-queue/app/schema"
)

type OutService interface {
	List(ctx context.Context, query *schema.OutMsgQueryParam) (*[]models.OutMessage, *paging.Paging, error)
	Publish(ctx context.Context, message *models.OutMessage) error
	CronResend() error
}
