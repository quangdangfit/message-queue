package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/quangdangfit/gosdk/utils/logger"
	"gopkg.in/mgo.v2/bson"

	dbs "gomq/app/database"
	"gomq/app/models"
)

type inMessageRepo struct {
	db dbs.IDatabase
}

func NewInMessageRepository(db dbs.IDatabase) InMessageRepository {
	return &inMessageRepo{db: db}
}

func (i *inMessageRepo) GetSingleInMessage(query map[string]interface{}) (*models.InMessage, error) {
	message := models.InMessage{}
	err := i.db.FindOne(dbs.CollectionInMessage, query, "-_id", &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (i *inMessageRepo) GetInMessages(query map[string]interface{}, limit int) (*[]models.InMessage, error) {
	message := []models.InMessage{}

	_, err := i.db.FindManyPaging(dbs.CollectionInMessage, query, "-_id", 1,
		limit, &message)
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

	err := i.db.InsertOne(dbs.CollectionInMessage, message)
	if err != nil {
		return err
	}
	return nil
}

func (i *inMessageRepo) UpdateInMessage(message *models.InMessage) error {
	selector := bson.M{"id": message.ID}

	var payload map[string]interface{}
	message.UpdatedTime = time.Now()
	data, _ := bson.Marshal(message)
	bson.Unmarshal(data, &payload)

	change := bson.M{"$set": payload}
	err := i.db.UpdateOne(dbs.CollectionInMessage, selector, change)
	if err != nil {
		return err
	}

	return nil
}

func (i *inMessageRepo) UpsertInMessage(message *models.InMessage) error {
	msg, err := i.GetSingleInMessage(bson.M{"id": message.ID})
	if msg != nil {
		err = i.UpdateInMessage(message)
		if err != nil {
			logger.Errorf("Failed to update msg %s, %s", message.ID, err)
			return err
		}

		logger.Infof("Updated msg %s", message.ID)
		return nil
	}

	err = i.AddInMessage(message)
	if err != nil {
		logger.Errorf("Failed to insert msg %s, %s", message.ID, err)
		return err
	}
	logger.Infof("Inserted msg %s", message.ID)
	return nil
}
