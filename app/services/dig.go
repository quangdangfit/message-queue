package services

import "go.uber.org/dig"

func Inject(container *dig.Container) error {
	_ = container.Provide(NewInMessageService)
	_ = container.Provide(NewOutMessageService)

	return nil
}
