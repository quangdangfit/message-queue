package dbs

import (
	"gomq/config"
	db "transport/lib/database"

	"gopkg.in/mgo.v2"
)

var Database db.MongoDB

const (
	CollectionInMessage  = "in_message"
	CollectionOutMessage = "out_message"
	CollectionRoutingKey = "routing_key"
)

func init() {
	dbConfig := db.DBConfig{
		MongoDBHosts: config.Config.MongoDB.Host,
		AuthDatabase: "admin",
		AuthUserName: config.Config.MongoDB.Username,
		AuthPassword: config.Config.MongoDB.Password,
		Database:     config.Config.MongoDB.Database,
		Env:          config.Config.MongoDB.Env,
		Replica:      config.Config.MongoDB.Replica,
	}

	Database = db.NewConnection(dbConfig)

	ensureIndex()
}

func ensureIndex() {
	// Index for OutMessages
	index1 := mgo.Index{
		Key:        []string{"origin_id, origin_model"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
		Name:       "origin_object_index",
	}
	Database.DropIndex(CollectionOutMessage, index1.Name)
	Database.EnsureIndex(CollectionOutMessage, index1)

	// Index for InMessages
	index2 := mgo.Index{
		Key:        []string{"origin_id, origin_model"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
		Name:       "origin_object_index",
	}
	Database.DropIndex(CollectionInMessage, index2.Name)
	Database.EnsureIndex(CollectionInMessage, index2)
}
