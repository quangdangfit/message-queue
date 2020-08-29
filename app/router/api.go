package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"

	"message-queue/app/api"
)

func RegisterAPI(e *gin.Engine, container *dig.Container) error {
	err := container.Invoke(func(
		outMsg *api.OutMsg,
		routing *api.Routing,
	) error {
		apiRoute := e.Group("/api/v1")
		apiRoute.GET("/out_messages", outMsg.List)
		apiRoute.POST("/out_messages", outMsg.Publish)
		apiRoute.PUT("/out_messages/:id", outMsg.Update)

		// Routing Keys
		apiRoute.GET("/routing_keys", routing.List)
		apiRoute.POST("/routing_keys", routing.Create)
		apiRoute.GET("/routing_keys/:id", routing.Retrieve)
		apiRoute.PUT("/routing_keys/:id", routing.Update)

		return nil
	})

	return err
}
