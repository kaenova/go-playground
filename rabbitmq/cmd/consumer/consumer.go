package main

import (
	"log"

	"github.com/kaenova/go-playground/rabbitmq/pkg"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	ch := InitChannelConnection()

	// Consume
	dev, err := ch.Consume(pkg.TIMED_QUEUE_NAME, "", true, false, false, true, amqp.Table{})
	if err != nil {
		log.Fatal(err.Error())
	}
	go func() {
		log.Println("Listening...")
		for data := range dev {
			log.Println(string(data.Body))
		}
	}()
}

func InitChannelConnection() *amqp.Channel {
	// Create connection
	conn, err := amqp.Dial("amqp://admin:password@localhost:5672")
	if err != nil {
		log.Fatal(err.Error())
	}

	// Create multiplex channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err.Error())
	}

	return ch
}
