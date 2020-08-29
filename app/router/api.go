package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"

	"message-queue/app/api"
)

func RegisterAPI(e *gin.Engine, container *dig.Container) error {
	err := container.Invoke(func(
		sender *api.Sender,
		routing *api.Routing,
	) error {
		apiRoute := e.Group("/api/v1/queue")
		apiRoute.POST("/messages", sender.PublishMessage)

		apiRoute.POST("/routing_keys", routing.Create)

		return nil
	})

	return err
}
