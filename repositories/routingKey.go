package repositories

import (
	"gomq/dbs"
	"gomq/models"
	"gopkg.in/mgo.v2/bson"
)

type RoutingKeyRepository interface {
	GetRoutingKey(name string) (*models.RoutingKey, error)
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
