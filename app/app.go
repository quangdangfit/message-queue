package app

import (
	"github.com/quangdangfit/gosdk/utils/logger"
	"go.uber.org/dig"

	"message-queue/app/api"
	"message-queue/app/dbs"
	"message-queue/app/grpc"
	"message-queue/app/queue"
	repoImpl "message-queue/app/repositories/impl"
	serviceImpl "message-queue/app/services/impl"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	// Inject repositories
	err := dbs.Inject(container)
	if err != nil {
		logger.Error("Failed to inject database", err)
	}

	// Inject repositories
	err = repoImpl.Inject(container)
	if err != nil {
		logger.Error("Failed to inject repositories", err)
	}

	// Inject services
	err = serviceImpl.Inject(container)
	if err != nil {
		logger.Error("Failed to inject services", err)
	}

	// Inject queue
	err = queue.Inject(container)
	if err != nil {
		logger.Error("Failed to inject queue", err)
	}

	// Inject APIs
	err = api.Inject(container)
	if err != nil {
		logger.Error("Failed to inject APIs", err)
	}

	// Inject RPC
	err = grpc.Inject(container)
	if err != nil {
		logger.Error("Failed to inject RPC", err)
	}

	return container
}
