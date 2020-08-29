package impl

import (
	"encoding/json"

	"github.com/jinzhu/copier"
	"github.com/quangdangfit/gosdk/utils/paging"
	"gopkg.in/mgo.v2/bson"

	"message-queue/app/dbs"
	"message-queue/app/models"
	"message-queue/app/repositories"
	"message-queue/app/schema"
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

func (r *routing) Get(query map[string]interface{}) (*models.RoutingKey, error) {
	var routingKey models.RoutingKey
	query["active"] = true
	err := r.db.FindOne(models.CollectionRoutingKey, query, "",
		&routingKey)
	if err != nil {
		return nil, err
	}
	return &routingKey, nil
}

func (r *routing) GetPrevious(srcRouting models.RoutingKey) (
	*models.RoutingKey, error) {

	query := bson.M{
		"group": srcRouting.Group,
		"value": srcRouting.Value - 1,
	}
	return r.Get(query)
}

func (r *routing) List(query *schema.RoutingQueryParam) (*[]models.RoutingKey, *paging.Paging, error) {
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
