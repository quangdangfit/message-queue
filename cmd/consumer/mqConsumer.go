package main

import "gomq/msgQueue"

func main() {
	consumer := msgQueue.NewConsumer()
	consumer.RunConsumer(nil)
}
