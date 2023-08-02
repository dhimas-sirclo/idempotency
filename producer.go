package main

import (
	"context"
	"encoding/json"
	"idempotency/faker"
	"log"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

// func startProducer(ch *amqp.Channel, queueName string) {
// 	log.Println("starting producer")
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	dot := 0
// 	for {
// 		// n := rand.Intn(100)
// 		// if n >= 99 {
// 		dot = dot + 1
// 		mod := ((dot % 2) + 1)
// 		content := fmt.Sprintf("%v%v", mod, strings.Repeat(".", mod))
// 		body := model.Body{
// 			Content: content,
// 		}
// 		b, err := json.Marshal(body)
// 		if err != nil {
// 			log.Println("producer", err)
// 			continue
// 		}
// 		err = ch.PublishWithContext(ctx,
// 			"",        // exchange
// 			queueName, // routing key
// 			false,     // mandatory
// 			false,     // immediate
// 			amqp.Publishing{
// 				// DeliveryMode: amqp.Persistent,
// 				ContentType: "text/plain",
// 				Body:        b,
// 				MessageId:   uuid.NewString(),
// 			})
// 		failOnError(err, "Failed to publish a message")
// 		log.Printf(" [x] Sent %s\n", body)

// 		// Don't burn cpu for demo.
// 		time.Sleep(time.Second * 3)
// 	}
// }

func upsertProductProducer(ch *amqp.Channel, queueName string) {
	log.Println("starting upsert product producer")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for {
		log.Println("..")
		body := faker.Product()
		b, err := json.Marshal(body)
		if err != nil {
			// log.Println("producer", err)
			log.Println(err)
			continue
		}
		err = ch.PublishWithContext(ctx,
			"",        // exchange
			queueName, // routing key
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				// DeliveryMode: amqp.Persistent,
				ContentType: "application/json",
				Body:        b,
				MessageId:   uuid.NewString(),
			})
		failOnError(err, "Failed to publish a message")
		// log.Printf(" [x] Sent %s\n", body)

		// Don't burn cpu for demo.
		time.Sleep(time.Second * 3)
	}
}
