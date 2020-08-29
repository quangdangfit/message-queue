package impl

import "go.uber.org/dig"

func Inject(container *dig.Container) error {
	_ = container.Provide(NewInService)
	_ = container.Provide(NewOutService)

	return nil
}
