package router

import (
	"github.com/labstack/echo"
	"go.uber.org/dig"

	"gomq/app/api"
)

func RegisterCron(e *echo.Echo, container *dig.Container) error {
	err := container.Invoke(func(
		cron *api.Cron,
	) error {
		cronRoute := e.Group("/api/v1/cron")
		cronRoute.POST("/resend", cron.Resend)

		return nil
	})

	return err
}
