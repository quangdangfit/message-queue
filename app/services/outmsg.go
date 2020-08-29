package services

import (
	"context"

	"gomq/app/models"
)

type OutService interface {
	Publish(ctx context.Context, message *models.OutMessage) error
	CronResend(limit int) error
}
