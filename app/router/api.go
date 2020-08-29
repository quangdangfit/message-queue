package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"

	"message-queue/app/api"
)

func RegisterAPI(e *gin.Engine, container *dig.Container) error {
	err := container.Invoke(func(
		sender *api.Sender,
	) error {
		msgRoute := e.Group("/api/v1/queue")
		msgRoute.POST("/messages", sender.PublishMessage)

		return nil
	})

	return err
}
