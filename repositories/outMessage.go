package repositories

import (
	"gomq/dbs"
	"gomq/models"
	"time"

	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
)

type OutMessageRepository interface {
	GetSingleOutMessage(query bson.M) (*models.OutMessage, error)
	GetOutMessages(query bson.M, limit int) (*[]models.OutMessage, error)
	AddOutMessage(message *models.OutMessage) error
	UpdateOutMessage(message *models.OutMessage) error
}

type outMessageRepo struct{}

func NewOutMessageRepo() OutMessageRepository {
	return &outMessageRepo{}
}

func (msg *outMessageRepo) GetSingleOutMessage(query bson.M) (*models.OutMessage, error) {
	message := models.OutMessage{}
	err := dbs.Database.FindOne(dbs.CollectionOutMessage, query, "-_id", &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}
func (msg *outMessageRepo) GetOutMessages(query bson.M, limit int) (*[]models.OutMessage, error) {
	message := []models.OutMessage{}
	_, err := dbs.Database.FindManyPaging(dbs.CollectionOutMessage, query, "-_id", 1,
		limit, &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (msg *outMessageRepo) AddOutMessage(message *models.OutMessage) error {
	message.CreatedTime = time.Now()
	message.UpdatedTime = time.Now()
	message.ID = uuid.New().String()

	err := dbs.Database.InsertOne(dbs.CollectionOutMessage, message)
	if err != nil {
		return err
	}
	return nil
}

func (msg *outMessageRepo) UpdateOutMessage(message *models.OutMessage) error {
	selector := bson.M{"id": message.ID}
	payload := map[string]interface{}{
		"updated_time": time.Now(),
		"status":       message.Status,
		"logs":         message.Logs,
	}

	change := bson.M{"$set": payload}
	err := dbs.Database.UpdateOne(dbs.CollectionOutMessage, selector, change)
	if err != nil {
		return err
	}

	return nil
}
