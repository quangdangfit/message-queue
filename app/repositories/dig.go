package repositories

import "go.uber.org/dig"

func Inject(container *dig.Container) error {
	_ = container.Provide(NewInRepository)

	return nil
}
