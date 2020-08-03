package router

import (
	"github.com/labstack/echo"
	"go.uber.org/dig"

	"gomq/app/api"
)

func RegisterAPI(e *echo.Echo, container *dig.Container) error {
	err := container.Invoke(func(
		sender *api.Sender,
	) error {
		msgRoute := e.Group("/api/v1/queue")
		msgRoute.POST("/messages", sender.PublishMessage)

		cronRoute := e.Group("/api/v1/cron")
		cronRoute.POST("/resend", sender.PublishMessage)

		return nil
	})

	return err
}
