package services

import (
	"context"
	"net/http"

	"gomq/app/models"
)

type InMessageService interface {
	HandleMessage(message *models.InMessage, routingKey string) error
	CronRetry(limit int) error
	CronRetryPrevious(limit int) error
	storeMessage(message *models.InMessage) (err error)
	callAPI(message *models.InMessage) (*http.Response, error)
}

type OutMessageService interface {
	Publish(ctx context.Context, message *models.OutMessage) error
	CronResend(limit int) error
}
