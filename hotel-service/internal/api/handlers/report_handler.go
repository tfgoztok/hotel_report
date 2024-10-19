package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/tfgoztok/hotel-service/internal/messaging"
)

// ReportHandler handles report-related requests
type ReportHandler struct {
	rabbitMQ *messaging.RabbitMQ // RabbitMQ instance for message publishing
}

// NewReportHandler creates a new instance of ReportHandler
func NewReportHandler(rabbitMQ *messaging.RabbitMQ) *ReportHandler {
	return &ReportHandler{rabbitMQ: rabbitMQ}
}

// ReportRequest represents the structure of a report request
type ReportRequest struct {
	ID     uuid.UUID `json:"id"`     // Unique identifier for the report
	Status string    `json:"status"` // Current status of the report
}

// RequestReport handles the HTTP request to create a new report request
func (h *ReportHandler) RequestReport(w http.ResponseWriter, r *http.Request) {
	reportRequest := ReportRequest{
		ID:     uuid.New(), // Generate a new UUID for the report
		Status: "pending",  // Set initial status to pending
	}

	err := h.rabbitMQ.PublishReportRequest("report_requests", reportRequest) // Publish the report request
	if err != nil {
		http.Error(w, "Failed to request report", http.StatusInternalServerError) // Handle publishing error
		return
	}

	w.WriteHeader(http.StatusAccepted)       // Respond with 202 Accepted status
	json.NewEncoder(w).Encode(reportRequest) // Encode and send the report request as JSON
}
