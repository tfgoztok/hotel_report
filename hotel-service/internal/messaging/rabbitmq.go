package messaging

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

// RabbitMQInterface defines the methods that our RabbitMQ implementation should have
type RabbitMQInterface interface {
	PublishReportRequest(queueName string, reportRequest interface{}) error
	Close()
}

// RabbitMQ struct implements the RabbitMQInterface
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewRabbitMQ creates a new RabbitMQ instance
func NewRabbitMQ(url string) (RabbitMQInterface, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	return &RabbitMQ{conn: conn, channel: ch}, nil
}

// PublishReportRequest publishes a report request to the specified queue
func (r *RabbitMQ) PublishReportRequest(queueName string, reportRequest interface{}) error {
	q, err := r.channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %v", err)
	}

	body, err := json.Marshal(reportRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal report request: %v", err)
	}

	err = r.channel.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %v", err)
	}

	return nil
}

// Close closes the RabbitMQ connection and channel
func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}
