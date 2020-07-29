package router

import (
	"github.com/labstack/echo"

	"gomq/packages/services"
)

func Message(e *echo.Echo) {
	msgRoute := e.Group("/api/v1/queue")

	sender := services.NewSender()
	msgRoute.POST("/messages", sender.PublishMessage)
}
