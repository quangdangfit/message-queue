package inrouting

import (
	"gomq/dbs"
	"gopkg.in/mgo.v2/bson"
)

type Repository interface {
	GetRoutingKey(query map[string]interface{}) (*RoutingKey, error)
	GetPreviousRoutingKey(srcRouting RoutingKey) (*RoutingKey, error)
}

type repo struct{}

func NewRepository() Repository {
	return &repo{}
}

func (r *repo) GetRoutingKey(query map[string]interface{}) (*RoutingKey, error) {
	var routingKey RoutingKey
	query["active"] = true
	err := dbs.Database.FindOne(dbs.CollectionRoutingKey, query, "",
		&routingKey)
	if err != nil {
		return nil, err
	}
	return &routingKey, nil
}

func (r *repo) GetPreviousRoutingKey(srcRouting RoutingKey) (
	*RoutingKey, error) {

	query := bson.M{
		"group": srcRouting.Group,
		"value": srcRouting.Value - 1,
	}
	return r.GetRoutingKey(query)
}
