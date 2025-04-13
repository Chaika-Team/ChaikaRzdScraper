package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

type Consumer struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewConsumer(rabbitMQURL string) (*Consumer, error) {
	log.Println("Connecting to RabbitMQ...")
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %v", err)
		conn.Close()
		return nil, err
	}

	log.Println("Connected to RabbitMQ and channel opened")

	return &Consumer{conn: conn, ch: ch}, nil
}

func (c *Consumer) Consume(queueName string) (<-chan amqp.Delivery, error) {
	return c.ch.Consume(
		queueName, // имя очереди
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
}

func (c *Consumer) Close() {
	log.Println("Closing RabbitMQ channel and connection...")
	if err := c.ch.Close(); err != nil {
		log.Printf("Failed to close channel: %v", err)
	}
	if err := c.conn.Close(); err != nil {
		log.Printf("Failed to close connection: %v", err)
	}
}
