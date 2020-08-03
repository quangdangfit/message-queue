package app

import (
	"github.com/quangdangfit/gosdk/utils/logger"
	"go.uber.org/dig"

	"gomq/app/api"
	"gomq/app/dbs"
	"gomq/app/queue"
	"gomq/app/repositories"
	"gomq/app/services"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	// Inject repositories
	err := dbs.Inject(container)
	if err != nil {
		logger.Error("Failed to inject database", err)
	}

	// Inject repositories
	err = repositories.Inject(container)
	if err != nil {
		logger.Error("Failed to inject repositories", err)
	}

	// Inject services
	err = services.Inject(container)
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

	return container
}
