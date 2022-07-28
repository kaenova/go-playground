package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kaenova/go-playground/rabbitmq/pkg"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Create connection
	conn, err := amqp.Dial("amqp://admin:password@localhost:5672")
	if err != nil {
		log.Fatal(err.Error())
	}
	// defer conn.Close()

	// Create multiplex channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err.Error())
	}
	// defer ch.Close()

	// Create queue
	q, err := ch.QueueDeclare(
		pkg.QUEUE_NAME,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx := context.Background()
	counter := 0
	for {
		time.Sleep(1 * time.Second)
		body := "hello world" + fmt.Sprint(counter)
		err = ch.PublishWithContext(
			ctx,
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			},
		)
		if err != nil {
			log.Fatal(err.Error())
		}

		counter++
	}
}
