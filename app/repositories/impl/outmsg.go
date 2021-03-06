package impl

import (
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/quangdangfit/gosdk/utils/paging"
	"gopkg.in/mgo.v2/bson"

	"message-queue/app/dbs"
	"message-queue/app/models"
	"message-queue/app/repositories"
	"message-queue/app/schema"
	"message-queue/config"
)

type outRepo struct {
	db dbs.IDatabase
}

func NewOutRepository(db dbs.IDatabase) repositories.OutRepository {
	return &outRepo{db: db}
}

func (o *outRepo) Retrieve(id string) (*models.OutMessage, error) {
	message := models.OutMessage{}
	query := bson.M{"id": id}
	err := o.db.FindOne(models.CollectionOutMessage, query, "-_id", &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (o *outRepo) Get(query *schema.OutMsgQueryParam) (*models.OutMessage, error) {
	message := models.OutMessage{}

	var mapQuery map[string]interface{}
	data, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &mapQuery)

	err = o.db.FindOne(models.CollectionOutMessage, mapQuery, "-_id", &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (o *outRepo) List(query *schema.OutMsgQueryParam) (*[]models.OutMessage, *paging.Paging, error) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = config.Config.PageLimit
	}

	var message []models.OutMessage
	var mapQuery map[string]interface{}
	data, err := json.Marshal(query)
	if err != nil {
		return nil, nil, err
	}
	json.Unmarshal(data, &mapQuery)

	pageInfo, err := o.db.FindManyPaging(models.CollectionOutMessage, mapQuery, "-_id", query.Page, query.Limit, &message)
	if err != nil {
		return nil, nil, err
	}

	return &message, pageInfo, nil
}

func (o *outRepo) Create(message *models.OutMessage) error {
	message.CreatedTime = time.Now()
	message.UpdatedTime = time.Now()
	message.ID = uuid.New().String()

	err := o.db.InsertOne(models.CollectionOutMessage, message)
	if err != nil {
		return err
	}
	return nil
}

func (o *outRepo) Update(id string, body *schema.OutMsgUpdateParam) (*models.OutMessage, error) {
	msg, err := o.Retrieve(id)
	if err != nil {
		return nil, err
	} else if msg == nil {
		return nil, errors.New("not found out message")
	}

	var update models.OutMessage
	copier.Copy(&update, &msg)
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &update)

	if reflect.DeepEqual(*msg, update) {
		return msg, nil
	}

	var value map[string]interface{}
	data, err = json.Marshal(update)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &value)

	selector := bson.M{"id": msg.ID}
	err = o.db.UpdateOne(models.CollectionOutMessage, selector, value)
	if err != nil {
		return nil, err
	}
	return &update, nil
}
