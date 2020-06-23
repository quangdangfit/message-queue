package dbs

import (
	"github.com/google/uuid"
	"gomq/models"
	"gopkg.in/mgo.v2/bson"
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

//TODO move to repositories package
func GetInMessageByStatus(status string, limit int) ([]models.InMessage, error) {
	message := []models.InMessage{}
	query := bson.M{"status": status}

	_, err := Database.FindManyPaging(CollectionInMessage, query, "-_id", 1,
		limit, &message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

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

func UpdateInMessage(message *models.InMessage) error {
	selector := bson.M{"id": message.ID}

	var payload map[string]interface{}
	message.UpdatedTime = time.Now()
	data, _ := bson.Marshal(message)
	bson.Unmarshal(data, &payload)

	change := bson.M{"$set": payload}
	err := Database.UpdateOne(CollectionInMessage, selector, change)
	if err != nil {
		return err
	}

	return nil
}
