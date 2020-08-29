package impl

import (
	"encoding/json"

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

func (r *routing) Retrieve(query map[string]interface{}) (*models.RoutingKey, error) {
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
	return r.Retrieve(query)
}

func (r *routing) List(query *schema.RoutingKeyQueryParam) (*[]models.RoutingKey, *paging.Paging, error) {
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
