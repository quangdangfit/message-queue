package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"

	"message-queue/app/api"
)

func RegisterAPI(e *gin.Engine, container *dig.Container) error {
	err := container.Invoke(func(
		outMsg *api.OutMsg,
		inMsg *api.InMsg,
		routing *api.Routing,
	) error {
		apiRoute := e.Group("/api/v1")

		// Out Messages
		apiRoute.GET("/out_messages", outMsg.List)
		apiRoute.POST("/out_messages", outMsg.Publish)
		apiRoute.PUT("/out_messages/:id", outMsg.Update)

		// In Messages
		apiRoute.GET("/in_messages", inMsg.List)

		// Routing Keys
		apiRoute.GET("/routing_keys", routing.List)
		apiRoute.POST("/routing_keys", routing.Create)
		apiRoute.GET("/routing_keys/:id", routing.Retrieve)
		apiRoute.PUT("/routing_keys/:id", routing.Update)

		return nil
	})

	return err
}
