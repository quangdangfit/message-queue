package grpc

import "go.uber.org/dig"

func Inject(container *dig.Container) error {
	_ = container.Provide(NewOutRPC)

	return nil
}
