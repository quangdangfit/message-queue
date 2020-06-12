package dbs

import (
	"gomq/models"
	"gopkg.in/mgo.v2/bson"
)

func GetRoutingKey(name string) (*models.RoutingKey, error) {
	var routingKey models.RoutingKey
	err := Database.FindOne(CollectionRoutingKey, bson.M{"name": name}, "", &routingKey)
	if err != nil {
		return nil, err
	}
	return &routingKey, nil
}
