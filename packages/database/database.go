package database

import (
	"gitlab.com/quangdangfit/gocommon/database"
	"gitlab.com/quangdangfit/gocommon/database/mongo"
	"go.uber.org/dig"
	"gopkg.in/mgo.v2"

	"gomq/config"
)

var Database mongo.MongoDB

const (
	CollectionInMessage  = "in_message"
	CollectionOutMessage = "out_message"
	CollectionRoutingKey = "routing_key"

	OutMessageStatusWait     = "wait"
	OutMessageStatusSent     = "sent"
	OutMessageStatusSentWait = "sent_wait"
	OutMessageStatusFailed   = "failed"
	OutMessageStatusCanceled = "canceled"
	OutMessageStatusInvalid  = "invalid"

	InMessageStatusReceived    = "received"
	InMessageStatusSuccess     = "success"
	InMessageStatusWaitRetry   = "wait_retry"
	InMessageStatusWorking     = "working"
	InMessageStatusFailed      = "failed"
	InMessageStatusInvalid     = "invalid"
	InMessageStatusWaitPrevMsg = "wait_prev_msg"
	InMessageStatusCanceled    = "canceled"
)

func NewDatabase() database.Database {
	dbConfig := database.DBConfig{
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

func Inject(container *dig.Container) {
	_ = container.Provide(NewDatabase)
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
