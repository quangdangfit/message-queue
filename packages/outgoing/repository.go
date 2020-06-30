package outgoing

import (
	"gomq/dbs"
	"time"

	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
)

type OutMessageRepository interface {
	GetSingleOutMessage(query bson.M) (*OutMessage, error)
	GetOutMessages(query bson.M, limit int) (*[]OutMessage, error)
	AddOutMessage(message *OutMessage) error
	UpdateOutMessage(message *OutMessage) error
}

type outMessageRepo struct{}

func NewOutMessageRepo() OutMessageRepository {
	return &outMessageRepo{}
}

func (msg *outMessageRepo) GetSingleOutMessage(query bson.M) (*OutMessage, error) {
	message := OutMessage{}
	err := dbs.Database.FindOne(dbs.CollectionOutMessage, query, "-_id", &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}
func (msg *outMessageRepo) GetOutMessages(query bson.M, limit int) (*[]OutMessage, error) {
	message := []OutMessage{}
	_, err := dbs.Database.FindManyPaging(dbs.CollectionOutMessage, query, "-_id", 1,
		limit, &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (msg *outMessageRepo) AddOutMessage(message *OutMessage) error {
	message.CreatedTime = time.Now()
	message.UpdatedTime = time.Now()
	message.ID = uuid.New().String()

	err := dbs.Database.InsertOne(dbs.CollectionOutMessage, message)
	if err != nil {
		return err
	}
	return nil
}

func (msg *outMessageRepo) UpdateOutMessage(message *OutMessage) error {
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
