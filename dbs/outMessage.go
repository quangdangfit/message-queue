package dbs

import (
	"github.com/google/uuid"
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

func GetOutMessageByStatus(status string, limit int) ([]models.OutMessage, error) {
	message := []models.OutMessage{}
	query := bson.M{"status": status}

	_, err := Database.FindManyPaging(CollectionOutMessage, query, "-_id", 1,
		limit, &message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func AddOutMessage(message *models.OutMessage) (*models.OutMessage, error) {
	message.CreatedTime = time.Now()
	message.UpdatedTime = time.Now()
	message.ID = uuid.New().String()

	err := Database.InsertOne(CollectionOutMessage, message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func UpdateOutMessage(message *models.OutMessage) error {
	selector := bson.M{"id": message.ID}
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
