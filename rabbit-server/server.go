package main

import (
	"log"
	traitement "somei/rabbit-server/traitement"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:25672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"encaissemeent", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			err := traitement.TraitementSuperComplique(d.Body)
			if err != nil {
				log.Println(err)
				d.Acknowledger.Nack(d.DeliveryTag, false, true)
			}
			if err == nil {
				log.Println("demande traitée sans erreur")
				d.Acknowledger.Ack(d.DeliveryTag, true)
			}
		}
	}()

	log.Printf(" [*] Serveur en ecout")
	<-forever
}
