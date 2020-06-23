package routers

import (
	"gomq/services"

	"github.com/labstack/echo"
)

func Message(e *echo.Echo) {
	msgRoute := e.Group("/api/v1/queue/")

	sender := services.NewSender()
	msgRoute.POST("messages", sender.PublishMessage)
}
