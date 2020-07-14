package outgoing

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"

	"gomq/dbs"
)

type Repository interface {
	GetSingleOutMessage(query map[string]interface{}) (*OutMessage, error)
	GetOutMessages(query map[string]interface{}, limit int) (*[]OutMessage, error)
	AddOutMessage(message *OutMessage) error
	UpdateOutMessage(message *OutMessage) error
}

type outRepo struct{}

func NewRepository() Repository {
	return &outRepo{}
}

func (repo *outRepo) GetSingleOutMessage(query map[string]interface{}) (*OutMessage, error) {
	message := OutMessage{}
	err := dbs.Database.FindOne(dbs.CollectionOutMessage, query, "-_id", &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}
func (repo *outRepo) GetOutMessages(query map[string]interface{}, limit int) (*[]OutMessage, error) {
	message := []OutMessage{}
	_, err := dbs.Database.FindManyPaging(dbs.CollectionOutMessage, query, "-_id", 1,
		limit, &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (repo *outRepo) AddOutMessage(message *OutMessage) error {
	message.CreatedTime = time.Now()
	message.UpdatedTime = time.Now()
	message.ID = uuid.New().String()

	err := dbs.Database.InsertOne(dbs.CollectionOutMessage, message)
	if err != nil {
		return err
	}
	return nil
}

func (repo *outRepo) UpdateOutMessage(message *OutMessage) error {
	selector := bson.M{"id": message.ID}

	var payload map[string]interface{}
	message.UpdatedTime = time.Now()
	data, _ := bson.Marshal(message)
	bson.Unmarshal(data, &payload)

	change := bson.M{"$set": payload}
	err := dbs.Database.UpdateOne(dbs.CollectionOutMessage, selector, change)
	if err != nil {
		return err
	}

	return nil
}
