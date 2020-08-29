package services

import (
	"context"

	"message-queue/app/models"
)

type OutService interface {
	Publish(ctx context.Context, message *models.OutMessage) error
	CronResend(limit int) error
}
