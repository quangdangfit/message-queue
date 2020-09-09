package main

import (
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"

	"github.com/quangdangfit/gosdk/utils/logger"

	"message-queue/app"
	"message-queue/app/grpc"
	"message-queue/app/router"
	"message-queue/app/services"
	"message-queue/config"
)

func main() {
	// Build DIG container
	container := app.BuildContainer()

	//Init server
	e := router.Initialize(container)

	// Start by mode
	if config.Config.Mode == 0 || config.Config.Mode == 1 {
		go func() {
			port := "8080"
			logger.Info("Listening at port: " + port)
			err := e.Run(":" + port)
			if err != nil && err != http.ErrServerClosed {
				logger.Fatal(err)
			}
		}()
	}

	if config.Config.Mode == 0 || config.Config.Mode == 2 {
		container.Invoke(func(
			inService services.InService,
		) {
			go inService.Consume()
		})
	}

	// Run RPC for publishing
	container.Invoke(func(
		outRPC *grpc.OutRPC,
	) {
		rpc.RegisterName("OutRPC", outRPC)

		logger.Info("Listening RPC at port 1234!")
		listener, err := net.Listen("tcp", ":1234")
		if err != nil {
			logger.Fatal("ListenTCP error: ", err)
		}

		for {
			conn, err := listener.Accept()
			if err != nil {
				logger.Fatal("Accept error: ", err)
			}
			go rpc.ServeConn(conn)
		}
	})

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	close(quit)
	logger.Info("Shutting down")
}
