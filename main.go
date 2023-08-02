package main

import (
	"os"
	"os/signal"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	redis "github.com/redis/go-redis/v9"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"UPSERT_PRODUCT", // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	var m sync.RWMutex
	// go _startConsumer(&m, ch, q.Name, "One", redisClient)
	// go _startConsumer(&m, ch, q.Name, "Two", redisClient)
	// go _startConsumer(&m, ch, q.Name, "Three", redisClient)
	// go _startConsumer(&m, ch, q.Name, "Four", redisClient)
	// go _startConsumer(&m, ch, q.Name, "Five", redisClient)
	go upsertProductConsumer(&m, ch, q.Name, "Product One", redisClient)
	go upsertProductConsumer(&m, ch, q.Name, "Product Two", redisClient)

	// go startProducer(ch, q.Name)
	go upsertProductProducer(ch, q.Name)

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, os.Interrupt)

	<-wait
}
