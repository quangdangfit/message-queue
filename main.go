package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/quangdangfit/gosdk/utils/logger"

	"gomq/app"
	"gomq/app/queue"
	"gomq/app/router"
	"gomq/config"
)

func main() {
	// Build DIG container
	container := app.BuildContainer()

	//Init serv
	e := router.Initialize(container)

	// Start by mode
	if config.Config.Mode == 0 || config.Config.Mode == 1 {
		go func() {
			port := "8080"
			logger.Info("Starting at port: " + port)
			err := e.Start(":" + port)
			if err != nil {
				logger.Error(err)
			}
		}()
	}

	if config.Config.Mode == 0 || config.Config.Mode == 2 {
		container.Invoke(func(
			consumer queue.Consumer,
		) {
			go consumer.RunConsumer(nil)
		})
	}

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	close(quit)
	logger.Info("Shutting down")

}
