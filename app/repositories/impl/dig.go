package impl

import "go.uber.org/dig"

func Inject(container *dig.Container) error {
	_ = container.Provide(NewInRepository)
	_ = container.Provide(NewOutRepository)
	_ = container.Provide(NewRoutingRepository)

	return nil
}
