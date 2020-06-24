package repositories

import (
	"gomq/dbs"
	"gomq/models"
	"time"

	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
)

type InMessageRepository interface {
	GetSingleInMessage(query bson.M) (*models.InMessage, error)
	GetInMessages(query bson.M, limit int) (*[]models.InMessage, error)
	AddInMessage(message *models.InMessage) error
	UpdateInMessage(message *models.InMessage) error
}

type inMessageRepo struct{}

func NewInMessageRepo() InMessageRepository {
	return &inMessageRepo{}
}

func (msg *inMessageRepo) GetSingleInMessage(query bson.M) (*models.InMessage, error) {
	message := models.InMessage{}
	err := dbs.Database.FindOne(dbs.CollectionInMessage, query, "-_id", &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (msg *inMessageRepo) GetInMessages(query bson.M, limit int) (*[]models.InMessage, error) {
	message := []models.InMessage{}

	_, err := dbs.Database.FindManyPaging(dbs.CollectionInMessage, query, "-_id", 1,
		limit, &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (msg *inMessageRepo) AddInMessage(message *models.InMessage) error {
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

func (msg *inMessageRepo) UpdateInMessage(message *models.InMessage) error {
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
