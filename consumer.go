package main

import (
	"context"
	"encoding/json"
	"fmt"
	"idempotency/model"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	redis "github.com/redis/go-redis/v9"
)

// func startConsumer(m *sync.RWMutex, ch *amqp.Channel, queueName string, consumerName string, redisClient *redis.Client) {
// 	log.Printf("starting consumer %v", consumerName)

// 	msgs, err := ch.Consume(
// 		queueName, // queue
// 		"",        // consumer
// 		false,     // auto-ack
// 		false,     // exclusive
// 		false,     // no-local
// 		false,     // no-wait
// 		nil,       // args
// 	)
// 	failOnError(err, "Failed to register a consumer")

// 	var forever chan struct{}

// 	go func() {
// 		for delivery := range msgs {
// 			var body model.Body
// 			err := json.Unmarshal(delivery.Body, &body)
// 			if err != nil {
// 				log.Println(err)
// 				delivery.Ack(false)
// 				continue
// 			}
// 			m.RLock()
// 			val, err := redisClient.Get(context.Background(), string(delivery.Body)).Result()
// 			m.RUnlock()
// 			log.Print(val)
// 			if err != nil && err == redis.Nil {
// 				// log.Println(err)
// 				m.Lock()
// 				redisClient.Set(context.Background(), string(delivery.Body), string(delivery.Body), time.Second*20).Err()
// 				m.Unlock()
// 			}
// 			if err == nil && val == string(delivery.Body) {
// 				log.Printf("Consumer %s receive duplicate message: %s, skipped", consumerName, delivery.Body)
// 				delivery.Ack(false)
// 				continue
// 			}
// 			log.Printf("Consumer %v received a message: %s", consumerName, delivery.Body)
// 			dotCount := bytes.Count(delivery.Body, []byte("."))
// 			t := time.Duration(dotCount)
// 			time.Sleep(t * time.Second)
// 			log.Printf("Message %s Done on %s", delivery.Body, consumerName)

// 			m.Lock()
// 			redisClient.Del(context.Background(), string(delivery.Body)).Err()
// 			m.Unlock()

// 			delivery.Ack(false)
// 		}
// 	}()

// 	// log.Printf("[*] Waiting for messages. To exit press CTRL+C")
// 	<-forever
// }

// func _startConsumer(m *sync.RWMutex, ch *amqp.Channel, queueName string, consumerName string, redisClient *redis.Client) {
// 	log.Printf("starting consumer %v", consumerName)

// 	msgs, err := ch.Consume(
// 		queueName, // queue
// 		"",        // consumer
// 		false,     // auto-ack
// 		false,     // exclusive
// 		false,     // no-local
// 		false,     // no-wait
// 		nil,       // args
// 	)
// 	failOnError(err, "Failed to register a consumer")

// 	var forever chan struct{}

// 	idem := NewRedis[model.Body](redisClient)

// 	go func() {
// 		for delivery := range msgs {
// 			var body model.Body
// 			err := json.Unmarshal(delivery.Body, &body)
// 			if err != nil {
// 				log.Println(err)
// 				delivery.Ack(false)
// 				continue
// 			}

// 			tr, has, err := idem.Start(context.Background(), delivery.MessageId)
// 			log.Println(tr, has, err, body)
// 			if err != nil {
// 				log.Println(consumerName, err)
// 			}
// 			if has {
// 				log.Println("duplicate: processed", consumerName, tr, body)
// 				delivery.Ack(false)
// 				continue
// 			}

// 			log.Printf("Consumer %v received a message: %s", consumerName, body)
// 			dotCount := bytes.Count([]byte(body.Content), []byte("."))
// 			t := time.Duration(dotCount)
// 			time.Sleep(t * time.Second)
// 			log.Printf("Message %v Done on %s", body, consumerName)

// 			err = idem.Store(context.Background(), delivery.MessageId, body)
// 			if err != nil {
// 				log.Println(consumerName, err)
// 			}

// 			delivery.Ack(false)
// 		}
// 	}()

// 	// log.Printf("[*] Waiting for messages. To exit press CTRL+C")
// 	<-forever
// }

func upsertProductConsumer(m *sync.RWMutex, ch *amqp.Channel, queueName string, consumerName string, redisClient *redis.Client) {
	log.Printf("starting consumer %v", consumerName)

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	idem := NewRedis[model.UpsertProductPayload](redisClient)

	go func() {
		for delivery := range msgs {
			log.Println(".")
			var body model.UpsertProductPayload
			err := json.Unmarshal(delivery.Body, &body)
			if err != nil {
				log.Println(err)
				delivery.Ack(false)
				continue
			}

			_, has, err := idem.Start(context.Background(), body.ActivityID)
			// log.Println(tr, has, err, body)
			if err != nil {
				log.Println(consumerName, err)
			}
			if has {
				log.Println("duplicate")
				// log.Println("duplicate: processed", consumerName, tr, body)
				delivery.Ack(false)
				continue
			}

			// remove duplicates
			products := unique(body.Products)
			fmt.Println(products)

			// log.Printf("Consumer %v received a message: %v", consumerName, body)
			time.Sleep(3 * time.Second)
			// log.Printf("Message %v Done on %s", body, consumerName)

			err = idem.Store(context.Background(), body.ActivityID, body)
			if err != nil {
				// log.Println(consumerName, err)
				log.Println(err)
			}

			delivery.Ack(false)
		}
	}()

	// log.Printf("[*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
