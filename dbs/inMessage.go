package dbs

import (
	"encoding/json"
	"transport/mq_service/models"

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

func AddInMessage(payload map[string]interface{}) (*models.InMessage, error) {
	var message models.InMessage

	routingKey := models.RoutingKey{}
	Database.FindOne(CollectionRoutingKey, bson.M{"name": payload["routing_key"]}, "", &routingKey)

	bytesMsg, _ := json.Marshal(payload)
	json.Unmarshal(bytesMsg, &message)
	message.RoutingKey = routingKey

	err := Database.InsertOne(CollectionInMessage, message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}
