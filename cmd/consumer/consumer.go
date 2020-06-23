package main

import (
	"gomq/queue"
)

func main() {
	consumer := queue.NewConsumer()
	consumer.RunConsumer(nil)
}
