package repositories

import (
	"gopkg.in/mgo.v2/bson"

	dbs "gomq/app/database"
	"gomq/app/models"
)

type routingRepo struct {
	db dbs.IDatabase
}

func NewRoutingRepository(db dbs.IDatabase) RoutingRepository {
	return &routingRepo{db: db}
}

func (r *routingRepo) GetRoutingKey(query map[string]interface{}) (*models.RoutingKey, error) {
	var routingKey models.RoutingKey
	query["active"] = true
	err := r.db.FindOne(models.CollectionRoutingKey, query, "",
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
