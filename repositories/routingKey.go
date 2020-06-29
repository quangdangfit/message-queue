package repositories

import (
	"gomq/dbs"
	"gomq/models"
	"gopkg.in/mgo.v2/bson"
)

type RoutingKeyRepository interface {
	GetRoutingKey(name string) (*models.RoutingKey, error)
	GetPreviousRoutingKey(srcRouting models.RoutingKey) (
		*models.RoutingKey, error)
}

type routingKeyRepo struct{}

func NewRoutingKeyRepo() RoutingKeyRepository {
	return &routingKeyRepo{}
}

func (r *routingKeyRepo) GetRoutingKey(name string) (*models.RoutingKey, error) {
	var routingKey models.RoutingKey
	query := bson.M{"name": name, "active": true}
	err := dbs.Database.FindOne(dbs.CollectionRoutingKey, query, "",
		&routingKey)
	if err != nil {
		return nil, err
	}
	return &routingKey, nil
}

func (r *routingKeyRepo) GetPreviousRoutingKey(srcRouting models.RoutingKey) (
	*models.RoutingKey, error) {

	var routingKey models.RoutingKey
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
