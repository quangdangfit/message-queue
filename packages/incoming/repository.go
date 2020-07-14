package incoming

import (
	"time"

	"github.com/google/uuid"
	"gitlab.com/quangdangfit/gocommon/utils/logger"
	"gopkg.in/mgo.v2/bson"

	"gomq/dbs"
)

type Repository interface {
	GetSingleInMessage(query map[string]interface{}) (*InMessage, error)
	GetInMessages(query map[string]interface{}, limit int) (*[]InMessage, error)
	AddInMessage(message *InMessage) error
	UpdateInMessage(message *InMessage) error
	UpsertInMessage(message *InMessage) error
}

type inRepo struct{}

func NewRepository() Repository {
	return &inRepo{}
}

func (repo *inRepo) GetSingleInMessage(query map[string]interface{}) (*InMessage, error) {
	message := InMessage{}
	err := dbs.Database.FindOne(dbs.CollectionInMessage, query, "-_id", &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (repo *inRepo) GetInMessages(query map[string]interface{}, limit int) (*[]InMessage, error) {
	message := []InMessage{}

	_, err := dbs.Database.FindManyPaging(dbs.CollectionInMessage, query, "-_id", 1,
		limit, &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (repo *inRepo) AddInMessage(message *InMessage) error {
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

func (repo *inRepo) UpdateInMessage(message *InMessage) error {
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

func (repo *inRepo) UpsertInMessage(message *InMessage) error {
	msg, err := repo.GetSingleInMessage(bson.M{"id": message.ID})
	if msg != nil {
		err = repo.UpdateInMessage(message)
		if err != nil {
			logger.Errorf("Failed to update msg %s, %s", message.ID, err)
			return err
		}

		logger.Infof("Updated msg %s", message.ID)
		return nil
	}

	err = repo.AddInMessage(message)
	if err != nil {
		logger.Errorf("Failed to insert msg %s, %s", message.ID, err)
		return err
	}
	logger.Infof("Inserted msg %s", message.ID)
	return nil
}
