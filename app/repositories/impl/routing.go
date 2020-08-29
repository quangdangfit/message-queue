package impl

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/jinzhu/copier"
	"github.com/quangdangfit/gosdk/utils/paging"
	"gopkg.in/mgo.v2/bson"

	"message-queue/app/dbs"
	"message-queue/app/models"
	"message-queue/app/repositories"
	"message-queue/app/schema"
	"message-queue/config"
)

type routing struct {
	db dbs.IDatabase
}

func NewRoutingRepository(db dbs.IDatabase) repositories.RoutingRepository {
	return &routing{db: db}
}

func (r *routing) Retrieve(id string) (*models.RoutingKey, error) {
	var routingKey models.RoutingKey
	query := bson.M{"id": id}
	err := r.db.FindOne(models.CollectionRoutingKey, query, "", &routingKey)
	if err != nil {
		return nil, err
	}
	return &routingKey, nil
}

func (r *routing) Get(query *schema.RoutingQueryParam) (*models.RoutingKey, error) {
	var routingKey models.RoutingKey
	var mapQuery map[string]interface{}
	data, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &mapQuery)
	mapQuery["active"] = true

	err = r.db.FindOne(models.CollectionRoutingKey, mapQuery, "", &routingKey)
	if err != nil {
		return nil, err
	}
	return &routingKey, nil
}

func (r *routing) List(query *schema.RoutingQueryParam) (*[]models.RoutingKey, *paging.Paging, error) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = config.Config.PageLimit
	}

	var routingKeys []models.RoutingKey
	var mapQuery map[string]interface{}
	data, err := json.Marshal(query)
	if err != nil {
		return nil, nil, err
	}
	json.Unmarshal(data, &mapQuery)

	pageInfo, err := r.db.FindManyPaging(models.CollectionRoutingKey, mapQuery, "-_id", query.Page, query.Limit, &routingKeys)
	if err != nil {
		return nil, nil, err
	}
	return &routingKeys, pageInfo, nil
}

func (r *routing) Create(body *schema.RoutingCreateParam) (*models.RoutingKey, error) {
	var routing models.RoutingKey
	copier.Copy(&routing, &body)
	routing.BeforeCreate()
	routing.Active = true

	var value map[string]interface{}
	data, err := json.Marshal(routing)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &value)

	err = r.db.InsertOne(models.CollectionRoutingKey, value)
	if err != nil {
		return nil, err
	}
	return &routing, nil
}

func (r *routing) Update(id string, body *schema.RoutingUpdateParam) (*models.RoutingKey, error) {
	routing, err := r.Retrieve(id)
	if err != nil {
		return nil, err
	} else if routing == nil {
		return nil, errors.New("not found routing key")
	}

	var update models.RoutingKey
	copier.Copy(&update, &routing)
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &update)

	if reflect.DeepEqual(*routing, update) {
		return routing, nil
	}

	var value map[string]interface{}
	data, err = json.Marshal(update)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &value)

	selector := bson.M{"id": routing.ID}
	err = r.db.UpdateOne(models.CollectionRoutingKey, selector, value)
	if err != nil {
		return nil, err
	}
	return &update, nil
}
