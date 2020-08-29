package dbs

import (
	"github.com/quangdangfit/gosdk/database"
	"github.com/quangdangfit/gosdk/database/mongo"
	"gopkg.in/mgo.v2"

	"message-queue/app/models"
	"message-queue/config"
)

var Database database.Mongo

type IDatabase interface {
	database.Mongo
}

func NewDatabase() IDatabase {
	dbConfig := database.Config{
		Hosts:        config.Config.MongoDB.Host,
		AuthDatabase: "admin",
		AuthUserName: config.Config.MongoDB.Username,
		AuthPassword: config.Config.MongoDB.Password,
		Database:     config.Config.MongoDB.Database,
		Env:          config.Config.MongoDB.Env,
		Replica:      config.Config.MongoDB.Replica,
	}

	return mongo.New(dbConfig)
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
	Database.DropIndex(models.CollectionOutMessage, indexOutId.Name)
	Database.EnsureIndex(models.CollectionOutMessage, indexOutId)

	indexOutObj := mgo.Index{
		Key:        []string{"origin_id, origin_model"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
		Name:       "origin_object_index",
	}
	Database.DropIndex(models.CollectionOutMessage, indexOutObj.Name)
	Database.EnsureIndex(models.CollectionOutMessage, indexOutObj)

	// Index for InMessages
	indexInId := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   false,
		Background: true,
		Sparse:     true,
		Name:       "in_message_id_index",
	}
	Database.DropIndex(models.CollectionInMessage, indexInId.Name)
	Database.EnsureIndex(models.CollectionInMessage, indexInId)

	indexInObj := mgo.Index{
		Key:        []string{"origin_id, origin_model"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
		Name:       "origin_object_index",
	}
	Database.DropIndex(models.CollectionInMessage, indexInObj.Name)
	Database.EnsureIndex(models.CollectionInMessage, indexInObj)
}
