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
	indexOutId := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   false,
		Background: true,
		Sparse:     true,
		Name:       "out_message_id_index",
	}
	Database.DropIndex(CollectionOutMessage, indexOutId.Name)
	Database.EnsureIndex(CollectionOutMessage, indexOutId)

	indexOutObj := mgo.Index{
		Key:        []string{"origin_id, origin_model"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
		Name:       "origin_object_index",
	}
	Database.DropIndex(CollectionOutMessage, indexOutObj.Name)
	Database.EnsureIndex(CollectionOutMessage, indexOutObj)

	// Index for InMessages
	indexInId := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   false,
		Background: true,
		Sparse:     true,
		Name:       "in_message_id_index",
	}
	Database.DropIndex(CollectionInMessage, indexInId.Name)
	Database.EnsureIndex(CollectionInMessage, indexInId)

	indexInObj := mgo.Index{
		Key:        []string{"origin_id, origin_model"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
		Name:       "origin_object_index",
	}
	Database.DropIndex(CollectionInMessage, indexInObj.Name)
	Database.EnsureIndex(CollectionInMessage, indexInObj)
}
