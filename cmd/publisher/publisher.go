package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/quangdangfit/gosdk/utils/logger"

	"gomq/app"
	"gomq/app/router"
)

func main() {
	container := app.BuildContainer()
	e := router.Initialize(container)

	// Start server
	go func() {
		port := "8080"
		logger.Info("Starting at port: " + port)
		err := e.Start(":" + port)
		if err != nil {
			logger.Error(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
