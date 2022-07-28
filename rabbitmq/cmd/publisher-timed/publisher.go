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
	var err error
	ch := InitChannelConnection()

	SetupExchangeAndQueue(ch)

	// Publish
	ctx := context.Background()
	counter := 0
	for {
		time.Sleep(1 * time.Second)
		body := "hello world" + fmt.Sprint(counter)
		log.Println("Sending", body)
		err = ch.PublishWithContext(
			ctx,
			pkg.TIMED_EXCHANGE_NAME,
			"",
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
				Headers: amqp.Table{
					"x-delay": 5000,
				},
			},
		)
		if err != nil {
			log.Fatal(err.Error())
		}

		counter++
	}
}

func SetupExchangeAndQueue(ch *amqp.Channel) {
	// Create exchange
	err := ch.ExchangeDeclare(
		pkg.TIMED_EXCHANGE_NAME, "x-delayed-message",
		true, false, false, false,
		amqp.Table{
			"x-delayed-type": "direct",
		},
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Create queue
	_, err = ch.QueueDeclare(
		pkg.TIMED_QUEUE_NAME,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Bind Queue and Exchange together
	err = ch.QueueBind(pkg.TIMED_QUEUE_NAME, "", pkg.TIMED_EXCHANGE_NAME, false, amqp.Table{})
	if err != nil {
		log.Fatal(err.Error())
	}
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
