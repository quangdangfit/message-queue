package services

import (
	"context"
	"net/http"

	"gomq/app/models"
)

type InMessageService interface {
	Consume()
	CronRetry(limit int) error
	CronRetryPrevious(limit int) error
	handle(message *models.InMessage, routingKey string) error
	storeMessage(message *models.InMessage) (err error)
	callAPI(message *models.InMessage) (*http.Response, error)
}

type OutMessageService interface {
	Publish(ctx context.Context, message *models.OutMessage) error
	CronResend(limit int) error
}
