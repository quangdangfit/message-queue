package repositories

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"

	dbs "gomq/app/database"
	"gomq/app/models"
)

type outMessageRepo struct {
	db dbs.IDatabase
}

func NewOutMessageRepository(db dbs.IDatabase) OutMessageRepository {
	return &outMessageRepo{db: db}
}

func (o *outMessageRepo) GetSingleOutMessage(query map[string]interface{}) (*models.OutMessage, error) {
	message := models.OutMessage{}
	err := o.db.FindOne(dbs.CollectionOutMessage, query, "-_id", &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}
func (o *outMessageRepo) GetOutMessages(query map[string]interface{}, limit int) (*[]models.OutMessage, error) {
	message := []models.OutMessage{}
	_, err := o.db.FindManyPaging(dbs.CollectionOutMessage, query, "-_id", 1,
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

	err := o.db.InsertOne(dbs.CollectionOutMessage, message)
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
	err := o.db.UpdateOne(dbs.CollectionOutMessage, selector, change)
	if err != nil {
		return err
	}

	return nil
}
