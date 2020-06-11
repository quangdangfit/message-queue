package route

import (
	"github.com/labstack/echo"
	"gomq/service"
)

func Message(e *echo.Echo) {
	msgRoute := e.Group("/api/v1/mq_service/")
	msgRoute.POST("messages", service.PublishMessage)
}
