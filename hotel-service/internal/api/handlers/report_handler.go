package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/olivere/elastic/v7"
	"github.com/tfgoztok/hotel-service/internal/messaging"
)

// ReportHandler handles report-related requests
type ReportHandler struct {
	rabbitMQ messaging.RabbitMQInterface // Interface for RabbitMQ messaging
	esClient *elastic.Client              // Elasticsearch client
}

// NewReportHandler creates a new instance of ReportHandler
func NewReportHandler(rabbitMQ messaging.RabbitMQInterface, esClient *elastic.Client) *ReportHandler {
	return &ReportHandler{rabbitMQ: rabbitMQ, esClient: esClient}
}

// ReportRequest represents the structure of a report request
type ReportRequest struct {
    ID       uuid.UUID `json:"id"`       // Unique identifier for the report
    Status   string    `json:"status"`   // Status of the report request
    Location string    `json:"location"` // Location associated with the report
}

// RequestReport handles incoming report requests
func (h *ReportHandler) RequestReport(w http.ResponseWriter, r *http.Request) {
    var request ReportRequest
    // Decode the JSON request body into the ReportRequest struct
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest) // Return error if decoding fails
        return
    }

    request.ID = uuid.New() // Generate a new UUID for the report
    request.Status = "pending" // Set the initial status of the report

    // Publish the report request to the RabbitMQ queue
    err := h.rabbitMQ.PublishReportRequest("report_requests", request)
    if err != nil {
        http.Error(w, "Failed to request report", http.StatusInternalServerError) // Handle publishing error
        return
    }

    // Index the report request in Elasticsearch
    _, err = h.esClient.Index().
        Index("report_requests").
        Id(request.ID.String()).
        BodyJson(request).
        Do(r.Context())
    if err != nil {
        http.Error(w, "Failed to index report request", http.StatusInternalServerError) // Handle indexing error
        return
    }

    w.WriteHeader(http.StatusAccepted) // Respond with 202 Accepted status
    json.NewEncoder(w).Encode(request) // Encode the request as JSON and send it in the response
}
