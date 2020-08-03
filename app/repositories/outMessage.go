package repositories

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"

	"gomq/app/dbs"
	"gomq/app/models"
	"gomq/app/schema"
)

type outMessageRepo struct {
	db dbs.IDatabase
}

func NewOutMessageRepository(db dbs.IDatabase) OutMessageRepository {
	return &outMessageRepo{db: db}
}

func (o *outMessageRepo) GetOutMessageByID(id string) (*models.OutMessage, error) {
	message := models.OutMessage{}
	query := bson.M{"id": id}

	err := o.db.FindOne(models.CollectionOutMessage, query, "-_id", &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (o *outMessageRepo) GetSingleOutMessage(query *schema.OutMessageQueryParam) (*models.OutMessage, error) {
	message := models.OutMessage{}

	var mapQuery map[string]interface{}
	data, err := bson.Marshal(query)
	if err != nil {
		return nil, err
	}
	bson.Unmarshal(data, &mapQuery)

	err = o.db.FindOne(models.CollectionOutMessage, mapQuery, "-_id", &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}
func (o *outMessageRepo) GetOutMessages(query *schema.OutMessageQueryParam, limit int) (*[]models.OutMessage, error) {
	var message []models.OutMessage

	var mapQuery map[string]interface{}
	data, err := bson.Marshal(query)
	if err != nil {
		return nil, err
	}
	bson.Unmarshal(data, &mapQuery)

	_, err = o.db.FindManyPaging(models.CollectionOutMessage, mapQuery, "-_id", 1,
		limit, &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (o *outMessageRepo) AddOutMessage(message *models.OutMessage) error {
	message.CreatedTime = time.Now()
	message.UpdatedTime = time.Now()
	message.ID = uuid.New().String()

	err := o.db.InsertOne(models.CollectionOutMessage, message)
	if err != nil {
		return err
	}
	return nil
}

func (o *outMessageRepo) UpdateOutMessage(message *models.OutMessage) error {
	selector := bson.M{"id": message.ID}

	var payload map[string]interface{}
	message.UpdatedTime = time.Now()
	data, _ := bson.Marshal(message)
	bson.Unmarshal(data, &payload)

	change := bson.M{"$set": payload}
	err := o.db.UpdateOne(models.CollectionOutMessage, selector, change)
	if err != nil {
		return err
	}

	return nil
}
