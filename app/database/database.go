package database

import (
	godb "github.com/quangdangfit/gosdk/database"
	"github.com/quangdangfit/gosdk/database/mongo"
	"go.uber.org/dig"
	"gopkg.in/mgo.v2"

	"gomq/app/models"
	"gomq/config"
)

var Database godb.Mongo

type IDatabase interface {
	godb.Mongo
}

func NewDatabase() IDatabase {
	dbConfig := godb.Config{
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

func Inject(container *dig.Container) error {
	_ = container.Provide(NewDatabase)

	return nil
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
