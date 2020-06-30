package incomming

import (
	"gomq/dbs"
	"time"

	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
)

type Repository interface {
	GetSingleInMessage(query map[string]interface{}) (*InMessage, error)
	GetInMessages(query map[string]interface{}, limit int) (*[]InMessage, error)
	AddInMessage(message *InMessage) error
	UpdateInMessage(message *InMessage) error
}

type inRepo struct{}

func NewInMessageRepo() Repository {
	return &inRepo{}
}

func (msg *inRepo) GetSingleInMessage(query map[string]interface{}) (*InMessage, error) {
	message := InMessage{}
	err := dbs.Database.FindOne(dbs.CollectionInMessage, query, "-_id", &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (msg *inRepo) GetInMessages(query map[string]interface{}, limit int) (*[]InMessage, error) {
	message := []InMessage{}

	_, err := dbs.Database.FindManyPaging(dbs.CollectionInMessage, query, "-_id", 1,
		limit, &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (msg *inRepo) AddInMessage(message *InMessage) error {
	message.CreatedTime = time.Now()
	message.UpdatedTime = time.Now()
	message.ID = uuid.New().String()
	message.Attempts = 0

	err := dbs.Database.InsertOne(dbs.CollectionInMessage, message)
	if err != nil {
		return err
	}
	return nil
}

func (msg *inRepo) UpdateInMessage(message *InMessage) error {
	selector := bson.M{"id": message.ID}

	var payload map[string]interface{}
	message.UpdatedTime = time.Now()
	data, _ := bson.Marshal(message)
	bson.Unmarshal(data, &payload)

	change := bson.M{"$set": payload}
	err := dbs.Database.UpdateOne(dbs.CollectionInMessage, selector, change)
	if err != nil {
		return err
	}

	return nil
}
