package main

import (
	"gomq/app"
	"gomq/app/queue"
	"gomq/app/services"
)

func main() {
	container := app.BuildContainer()

	container.Invoke(func(
		service services.InMessageService,
	) error {
		consumer := queue.NewConsumer(service)
		consumer.RunConsumer(nil)

		return nil
	})
}
