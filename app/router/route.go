package router

import (
	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gosdk/utils/logger"
	"go.uber.org/dig"
)

func Initialize(container *dig.Container) *gin.Engine {
	app := gin.New()
	err := RegisterAPI(app, container)
	if err != nil {
		logger.Error("Failed to register API: ", err)
	}

	err = RegisterCron(app, container)
	if err != nil {
		logger.Error("Failed to register Cron API: ", err)
	}

	RegisterDocs(app)

	return app
}
