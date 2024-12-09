package queue

import (
	"log"

	"github.com/streadway/amqp"
)

var (
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
)

// InitQueue initializes the RabbitMQ connection and declares a queue
func InitQueue() {
	var err error
	// Connect to RabbitMQ
	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	// Open a channel
	ch, err = conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	// Declare a queue
	q, err = ch.QueueDeclare("image_processing", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	log.Println("RabbitMQ initialized successfully!")
}

// PublishToQueue publishes a list of image URLs to the RabbitMQ queue
func PublishToQueue(imageURLs []string) error {
	for _, url := range imageURLs {
		err := ch.Publish("", q.Name, false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(url),
		})
		if err != nil {
			log.Printf("Failed to publish message: %v", err)
			return err
		}
	}
	log.Println("Image URLs published to queue")
	return nil
}

// CloseQueue closes the RabbitMQ connection and channel
func CloseQueue() {
	if ch != nil {
		ch.Close()
	}
	if conn != nil {
		conn.Close()
	}
	log.Println("RabbitMQ connection closed")
}
