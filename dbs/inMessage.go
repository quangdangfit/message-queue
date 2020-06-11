package dbs

import (
	"gomq/models"

	"gopkg.in/mgo.v2/bson"
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

func AddInMessage(message *models.InMessage, strRoutingKey string) (
	*models.InMessage, error) {
	routingKey := models.RoutingKey{}
	Database.FindOne(CollectionRoutingKey,
		bson.M{"name": strRoutingKey}, "", &routingKey)

	message.RoutingKey = routingKey

	err := Database.InsertOne(CollectionInMessage, message)
	if err != nil {
		return nil, err
	}
	return message, nil
}
