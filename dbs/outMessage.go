package dbs

import (
	"gomq/models"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	OutMessageStatusWait     = "wait"
	OutMessageStatusSent     = "sent"
	OutMessageStatusSentWait = "sent_wait"
	OutMessageStatusFailed   = "failed"
	OutMessageStatusCancel   = "canceled"
	OutMessageStatusInvalid  = "invalid"
)

func GetWaitOutMessage(limit int) ([]models.OutMessage, error) {
	message := []models.OutMessage{}
	query := bson.M{"status": "wait"}

	_, err := Database.FindManyPaging(CollectionOutMessage, query, "-_id", 1,
		limit, &message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func AddOutMessage(message *models.OutMessage) (*models.OutMessage, error) {
	err := Database.InsertOne(CollectionOutMessage, message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func UpdateOutMessage(message *models.OutMessage) error {
	//TODO: create field id (uuid) and query by id
	selector := bson.M{
		"routing_key":  message.RoutingKey,
		"origin_model": message.OriginModel,
		"routing_code": message.OriginCode,
		"create_time":  message.CreatedTime,
	}

	payload := map[string]interface{}{
		"updated_time": time.Now(),
		"status":       message.Status,
		"logs":         message.Logs,
	}

	change := bson.M{"$set": payload}
	err := Database.UpdateOne(CollectionOutMessage, selector, change)
	if err != nil {
		return err
	}

	return nil
}
