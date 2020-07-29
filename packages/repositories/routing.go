package repositories

import (
	"gitlab.com/quangdangfit/gocommon/database"
	"gopkg.in/mgo.v2/bson"

	"gomq/dbs"
	"gomq/packages/models"
)

type routingRepo struct {
	db database.Database
}

func NewRoutingRepository(db database.Database) RoutingRepository {
	return &routingRepo{db: db}
}

func (r *routingRepo) GetRoutingKey(query map[string]interface{}) (*models.RoutingKey, error) {
	var routingKey models.RoutingKey
	query["active"] = true
	err := dbs.Database.FindOne(dbs.CollectionRoutingKey, query, "",
		&routingKey)
	if err != nil {
		return nil, err
	}
	return &routingKey, nil
}

func (r *routingRepo) GetPreviousRoutingKey(srcRouting models.RoutingKey) (
	*models.RoutingKey, error) {

	query := bson.M{
		"group": srcRouting.Group,
		"value": srcRouting.Value - 1,
	}
	return r.GetRoutingKey(query)
}
