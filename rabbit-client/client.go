package main

import (
	"encoding/json"
	"fmt"
	"log"

	uuid "github.com/satori/go.uuid"
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

	// Enveloppe
	mnt := 100
	// selectionner facture
	fac := uuid.NewV4()
	log.Println("Selection de la facture :" + fac.String())
	log.Println("Montant encaissé :" + fmt.Sprint(mnt))
	// Generation d'un numero de piece
	piece := uuid.NewV4()
	log.Println("Numero de la piece :" + piece.String())

	encaissement := map[string]string{
		"NumPiece": piece.String(),
		"NumFac":   fac.String(),
		"Mnt":      fmt.Sprint(mnt),
		"Caissier": "somei",
	}

	encaissementJson, err := json.Marshal(encaissement)

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         encaissementJson,
			DeliveryMode: amqp.Persistent,
		})

	failOnError(err, "[ ] Impossible d'executer l'encaissement essayez ulterieurement")
	log.Printf(" [x] Impression de %s", encaissement["NumPiece"])
}
