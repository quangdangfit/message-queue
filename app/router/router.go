package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/quangdangfit/gosdk/utils/logger"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/dig"

	_ "gomq/docs"
)

func Initialize(container *dig.Container) *echo.Echo {
	app := echo.New()
	app.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
	}))

	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339_nano}\t${method}\t${latency_human}\t${uri}\t${status}\n",
	}))
	app.Use(middleware.CORS())
	app.Use(middleware.RequestID())
	app.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	app.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}
		if code == http.StatusUnauthorized {
			_ = c.JSON(http.StatusUnauthorized, nil)
		}
		app.DefaultHTTPErrorHandler(err, c)
	}

	err := RegisterAPI(app, container)
	if err != nil {
		logger.Error("Failed to register API: ", err)
	}

	err = RegisterCron(app, container)
	if err != nil {
		logger.Error("Failed to register Cron API: ", err)
	}

	app.GET("/swagger/*", echoSwagger.WrapHandler)

	return app
}
