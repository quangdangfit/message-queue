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
		apiRoute := e.Group("/api/v1")
		apiRoute.POST("/out_messages", sender.PublishMessage)

		// Routing Keys
		apiRoute.GET("/routing_keys", routing.List)
		apiRoute.POST("/routing_keys", routing.Create)
		apiRoute.GET("/routing_keys/:id", routing.Retrieve)
		apiRoute.PUT("/routing_keys/:id", routing.Update)

		return nil
	})

	return err
}
