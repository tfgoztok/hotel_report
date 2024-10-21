// File: tests/unit/rabbitmq_test.go

package unit

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock RabbitMQ for unit tests
type MockRabbitMQ struct {
	publishedMessages [][]byte
}

func (m *MockRabbitMQ) PublishReportRequest(queueName string, reportRequest interface{}) error {
	body, _ := json.Marshal(reportRequest)
	m.publishedMessages = append(m.publishedMessages, body)
	return nil
}

func (m *MockRabbitMQ) Close() {}

// Unit test for PublishReportRequest
func TestPublishReportRequest(t *testing.T) {
	mock := &MockRabbitMQ{}

	reportRequest := struct {
		ID string `json:"id"`
	}{
		ID: "test-id",
	}

	err := mock.PublishReportRequest("test-queue", reportRequest)
	assert.NoError(t, err)
	assert.Len(t, mock.publishedMessages, 1)

	var receivedRequest struct {
		ID string `json:"id"`
	}
	err = json.Unmarshal(mock.publishedMessages[0], &receivedRequest)
	assert.NoError(t, err)
	assert.Equal(t, "test-id", receivedRequest.ID)
}
