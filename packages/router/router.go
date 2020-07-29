package router

import (
	"github.com/labstack/echo"

	"gomq/packages/api"
)

func Message(e *echo.Echo) {
	msgRoute := e.Group("/api/v1/queue")

	sender := api.NewSender()
	msgRoute.POST("/messages", sender.PublishMessage)
}
