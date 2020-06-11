package dbs

import (
	"transport/mq_service/models"
)

const (
	OutMessageStatusWait     = "wait"
	OutMessageStatusSent     = "sent"
	OutMessageStatusSentWait = "sent_wait"
	OutMessageStatusFailed   = "failed"
	OutMessageStatusCancel   = "canceled"
	OutMessageStatusInvalid  = "invalid"
)

func AddOutMessage(message *models.OutMessage) (*models.OutMessage, error) {
	err := Database.InsertOne(CollectionOutMessage, message)
	if err != nil {
		return nil, err
	}
	return message, nil
}
