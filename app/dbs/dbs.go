package dbs

import (
	"github.com/quangdangfit/gosdk/database"
	"github.com/quangdangfit/gosdk/database/mongo"

	"message-queue/config"
)

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
