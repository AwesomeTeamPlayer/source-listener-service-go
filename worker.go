package main

import (
	"fmt"
	"os"
	"github.com/streadway/amqp"
	"log"
	"encoding/json"
)

var conn *amqp.Connection

type ReceivedEvent struct {
	ObjectType 	string 	`json:"objectType"`
	ObjectId	string	`json:"objectId"`
}

func publishEvent(clientId *string, receivedEvent string) {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	var body string = "{\"clientId\":\"" + *clientId + "\", " + receivedEvent[1:]

	err = ch.Publish(
		os.Getenv("AMQP_EXCAHNGE_NAME"),
		"",
		false,
		false,
		amqp.Publishing {
			ContentType: "application/json",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	fmt.Println("Event published: ", body)
}

func runWorker() {
	fmt.Println("Start worker")

	var connectString string = "amqp://" + os.Getenv("AMQP_USER") + ":" + os.Getenv("AMQP_PASSWORD") + "@" + os.Getenv("AMQP_HOST") + ":" + os.Getenv("AMQP_PORT") + "/"

	fmt.Println("Try connect to Rabbit: " + connectString)

	var err error;
	conn, err = amqp.Dial(connectString)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	msgs, err := ch.Consume(
		os.Getenv("AMQP_QUEUE"), // queue
		"",     // consumer
		false,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			event := ReceivedEvent{}
			json.Unmarshal([]byte(d.Body), &event)

			clientsIds := getClients(event.ObjectType, event.ObjectId)

			for _, clientId := range clientsIds {
				publishEvent(&clientId, string(d.Body))
			}

			d.Ack(true)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}