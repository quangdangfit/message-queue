package dbs

import (
	"github.com/google/uuid"
	"gomq/models"
	"time"
)

const (
	InMessageStatusReceived     = "received"
	InMessageStatusSuccess      = "success"
	InMessageStatusWaitRetry    = "wait_retry"
	InMessageStatusWorking      = "working"
	InMessageStatusFailed       = "failed"
	InMessageStatusInvalid      = "invalid"
	InMessageStatusWaitPrevMsg  = "wait_prev_msg"
	InMessageStatusWaitCanceled = "canceled"
)

func AddInMessage(message *models.InMessage) (*models.InMessage, error) {
	message.CreatedTime = time.Now()
	message.UpdatedTime = time.Now()
	message.ID = uuid.New().String()

	err := Database.InsertOne(CollectionInMessage, message)
	if err != nil {
		return nil, err
	}
	return message, nil
}
