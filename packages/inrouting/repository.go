package inrouting

import (
	"gomq/dbs"
	"gopkg.in/mgo.v2/bson"
)

type Repository interface {
	GetRoutingKey(name string) (*RoutingKey, error)
	GetPreviousRoutingKey(srcRouting RoutingKey) (*RoutingKey, error)
}

type repo struct{}

func NewRoutingKeyRepo() Repository {
	return &repo{}
}

func (r *repo) GetRoutingKey(name string) (*RoutingKey, error) {
	var routingKey RoutingKey
	query := bson.M{"name": name, "active": true}
	err := dbs.Database.FindOne(dbs.CollectionRoutingKey, query, "",
		&routingKey)
	if err != nil {
		return nil, err
	}
	return &routingKey, nil
}

func (r *repo) GetPreviousRoutingKey(srcRouting RoutingKey) (
	*RoutingKey, error) {

	var routingKey RoutingKey
	query := bson.M{
		"group":  srcRouting.Group,
		"active": true,
		"value":  srcRouting.Value - 1,
	}

	err := dbs.Database.FindOne(dbs.CollectionRoutingKey, query, "",
		&routingKey)
	if err != nil {
		return nil, err
	}
	return &routingKey, nil
}
