package routers

import (
	"github.com/labstack/echo"
	"gomq/service"
)

func Message(e *echo.Echo) {
	msgRoute := e.Group("/api/v1/queue/")

	sender := service.NewSender()
	msgRoute.POST("messages", sender.PublishMessage)
}
