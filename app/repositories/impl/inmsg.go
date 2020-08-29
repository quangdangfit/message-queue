package impl

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"

	"message-queue/app/dbs"
	"message-queue/app/models"
	"message-queue/app/repositories"
	"message-queue/app/schema"
)

type inRepo struct {
	db dbs.IDatabase
}

func NewInRepository(db dbs.IDatabase) repositories.InRepository {
	return &inRepo{db: db}
}

func (i *inRepo) GetByID(id string) (*models.InMessage, error) {
	message := models.InMessage{}
	query := bson.M{"id": id}
	err := i.db.FindOne(models.CollectionInMessage, query, "-_id", &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (i *inRepo) Retrieve(query *schema.InMsgQueryParam) (*models.InMessage, error) {
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

func (i *inRepo) List(query *schema.InMsgQueryParam, limit int) (*[]models.InMessage, error) {
	var message []models.InMessage

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

func (i *inRepo) Create(message *models.InMessage) error {
	message.CreatedTime = time.Now()
	message.UpdatedTime = time.Now()
	message.ID = uuid.New().String()
	message.Attempts = 0

	var value map[string]interface{}
	data, err := json.Marshal(message)
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

func (i *inRepo) Update(message *models.InMessage) error {
	message.UpdatedTime = time.Now()
	selector := bson.M{"id": message.ID}

	var payload map[string]interface{}
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

func (i *inRepo) Upsert(message *models.InMessage) error {
	msg, _ := i.GetByID(message.ID)
	if msg != nil {
		return i.Update(message)
	}
	return i.Create(message)
}
