package repositories

import "go.uber.org/dig"

func Inject(container *dig.Container) error {
	_ = container.Provide(NewInMessageRepository)
	_ = container.Provide(NewOutMessageRepository)
	_ = container.Provide(NewRoutingRepository)

	return nil
}
