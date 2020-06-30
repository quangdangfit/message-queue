package main

import (
	"gomq/packages/queue"
)

func main() {
	consumer := queue.NewConsumer()
	consumer.RunConsumer(nil)
}
