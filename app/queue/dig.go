package queue

import (
	"go.uber.org/dig"
)

func Inject(container *dig.Container) error {
	_ = container.Provide(NewPublisher)
	_ = container.Provide(NewConsumer)
	return nil
}
