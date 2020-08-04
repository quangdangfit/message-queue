package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"

	"gomq/app/api"
)

func RegisterCron(e *gin.Engine, container *dig.Container) error {
	err := container.Invoke(func(
		cron *api.Cron,
	) error {
		cronRoute := e.Group("/api/v1/cron")
		cronRoute.POST("/resend", cron.Resend)
		cronRoute.POST("/retry", cron.Retry)
		cronRoute.POST("/retry_previous", cron.RetryPrevious)

		return nil
	})

	return err
}
