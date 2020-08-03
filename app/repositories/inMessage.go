package repositories

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"

	"gomq/app/dbs"
	"gomq/app/models"
	"gomq/app/schema"
)

type inMessageRepo struct {
	db dbs.IDatabase
}

func NewInMessageRepository(db dbs.IDatabase) InMessageRepository {
	return &inMessageRepo{db: db}
}

func (i *inMessageRepo) GetInMessageByID(id string) (*models.InMessage, error) {
	message := models.InMessage{}
	query := bson.M{"id": id}
	err := i.db.FindOne(models.CollectionInMessage, query, "-_id", &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (i *inMessageRepo) GetSingleInMessage(query *schema.InMessageQueryParam) (*models.InMessage, error) {
	message := models.InMessage{}

	var mapQuery map[string]interface{}
	data, err := bson.Marshal(query)
	if err != nil {
		return nil, err
	}
	bson.Unmarshal(data, &mapQuery)

	err = i.db.FindOne(models.CollectionInMessage, mapQuery, "-_id", &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (i *inMessageRepo) GetInMessages(query *schema.InMessageQueryParam, limit int) (*[]models.InMessage, error) {
	message := []models.InMessage{}

	var mapQuery map[string]interface{}
	data, err := bson.Marshal(query)
	if err != nil {
		return nil, err
	}
	bson.Unmarshal(data, &mapQuery)

	_, err = i.db.FindManyPaging(models.CollectionInMessage, mapQuery, "-_id", 1, limit, &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (i *inMessageRepo) AddInMessage(message *models.InMessage) error {
	message.CreatedTime = time.Now()
	message.UpdatedTime = time.Now()
	message.ID = uuid.New().String()
	message.Attempts = 0

	var value map[string]interface{}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	json.Unmarshal(data, &value)

	err = i.db.InsertOne(models.CollectionInMessage, value)
	if err != nil {
		return err
	}
	return nil
}

func (i *inMessageRepo) UpdateInMessage(message *models.InMessage) error {
	selector := bson.M{"id": message.ID}

	var payload map[string]interface{}
	message.UpdatedTime = time.Now()
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	json.Unmarshal(data, &payload)

	change := bson.M{"$set": payload}
	err = i.db.UpdateOne(models.CollectionInMessage, selector, change)
	if err != nil {
		return err
	}

	return nil
}

func (i *inMessageRepo) UpsertInMessage(message *models.InMessage) error {
	msg, err := i.GetInMessageByID(message.ID)
	if err != nil {
		return err
	}

	if msg != nil {
		return i.UpdateInMessage(message)
	}
	return i.AddInMessage(message)
}
