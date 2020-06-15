package dbs

import (
	"gomq/models"
	"gopkg.in/mgo.v2/bson"
)

func GetRoutingKey(name string) (*models.RoutingKey, error) {
	var routingKey models.RoutingKey
	query := bson.M{"name": name, "active": true}
	err := Database.FindOne(CollectionRoutingKey, query, "",
		&routingKey)
	if err != nil {
		return nil, err
	}
	return &routingKey, nil
}
