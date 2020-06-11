package msgHandler

import (
	"encoding/json"
	"fmt"
	"gomq/config"
	"gomq/dbs"
	"gomq/models"
	"gomq/msgQueue"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type Receiver interface {
	Consuming() error
	StoreMessage(byteMsg []byte) (*models.InMessage, error)
}

type receiver struct {
	consumer msgQueue.Consumer
}

func NewReceiver() Receiver {
	r := receiver{
		consumer: msgQueue.NewConsumer(),
	}
	return &r
}

func (r *receiver) Consuming() error {
	//msgChan, _ := r.consumer.Consume()
	//logger.Info(msgChan)
	//
	//return nil

	fmt.Println("Connecting to RabbitMQ")
	conn, err := amqp.Dial(config.Config.AMQP.URL) //Insert the  connection string
	failOnError(err, "Failed to connect to RabbitMQ server", "Successfully connected to RabbitMQ server")
	defer conn.Close()
	failOnError(err, "RabbitMQ connection failure", "RabbitMQ Connection Established")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel", "Opened the channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"truck.order", //name
		//"ha.monitoring",
		true,
		false, //delete when unused
		false, //exclusive
		false, //no-wait
		nil,   //arguements
	)
	failOnError(err, "Failed to declare the queue", "Declared the queue")

	msgs, err := ch.Consume(
		q.Name, //queue
		"",     //consumer
		true,   //auto-ack
		false,  //exclusive
		false,  //no-local
		false,  //no-wait
		nil,    //args
	)
	failOnError(err, "Failed to register a consumer ", "Registered the consumer")

	msgCount := 0
	go func() {
		for d := range msgs {

			msgCount++
			msg, _ := r.StoreMessage(d.Body)

			fmt.Printf("\nMessage Count: %d, Message Body: %s\n", msgCount, d.Body)
			fmt.Println(msg)

		}
	}()

	select {
	case <-time.After(time.Second * 50):
		fmt.Printf("Total Messages Fetched: %d\n", msgCount)
		fmt.Println("No more messages in queue. Timing out...")

	}

	return nil
}

func failOnError(err error, msgerr string, msgsuc string) {
	if err != nil {
		log.Fatalf("%s: %s", msgerr, err)

	} else {
		fmt.Printf("%s\n", msgsuc)
	}

}

func (r *receiver) StoreMessage(byteMsg []byte) (*models.InMessage, error) {
	msg := map[string]interface{}{}
	json.Unmarshal(byteMsg, &msg)

	message, _ := dbs.AddInMessage(msg)

	return message, nil
}
