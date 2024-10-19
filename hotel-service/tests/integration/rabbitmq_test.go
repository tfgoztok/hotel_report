package integration

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tfgoztok/hotel-service/internal/messaging"
)

func TestRabbitMQIntegration(t *testing.T) {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}

	// Connect to RabbitMQ
	rabbitMQ, err := messaging.NewRabbitMQ(rabbitMQURL)
	require.NoError(t, err)
	defer rabbitMQ.Close()

	// Create a separate connection for consuming messages
	conn, err := amqp.Dial(rabbitMQURL)
	require.NoError(t, err)
	defer conn.Close()

	ch, err := conn.Channel()
	require.NoError(t, err)
	defer ch.Close()

	// Declare a test queue
	queueName := "test_queue"
	_, err = ch.QueueDeclare(
		queueName,
		false, // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	require.NoError(t, err)

	// Publish a message
	testMessage := struct {
		ID string `json:"id"`
	}{
		ID: "integration-test-id",
	}
	err = rabbitMQ.PublishReportRequest(queueName, testMessage)
	require.NoError(t, err)

	// Consume the message
	msgs, err := ch.Consume(
		queueName,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	require.NoError(t, err)

	// Wait for the message
	select {
	case msg := <-msgs:
		var receivedMessage struct {
			ID string `json:"id"`
		}
		err = json.Unmarshal(msg.Body, &receivedMessage)
		assert.NoError(t, err)
		assert.Equal(t, "integration-test-id", receivedMessage.ID)
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for message")
	}

	// Clean up: delete the queue after the test
	_, err = ch.QueueDelete(queueName, false, false, false)
	require.NoError(t, err)
}

func TestNewRabbitMQ(t *testing.T) {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}

	rabbitMQ, err := messaging.NewRabbitMQ(rabbitMQURL)
	assert.NoError(t, err)
	assert.NotNil(t, rabbitMQ)

	// Test that we can publish a message
	err = rabbitMQ.PublishReportRequest("test_queue", struct{}{})
	assert.NoError(t, err)

	rabbitMQ.Close()

	// Clean up: delete the queue after the test
	conn, err := amqp.Dial(rabbitMQURL)
	require.NoError(t, err)
	defer conn.Close()

	ch, err := conn.Channel()
	require.NoError(t, err)
	defer ch.Close()

	_, err = ch.QueueDelete("test_queue", false, false, false)
	require.NoError(t, err)
}

func TestNewRabbitMQError(t *testing.T) {
	_, err := messaging.NewRabbitMQ("invalid_url")
	assert.Error(t, err)
}
