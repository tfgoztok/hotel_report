package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/olivere/elastic/v7"
	"github.com/tfgoztok/hotel-service/internal/messaging"
)

type ReportHandler struct {
	rabbitMQ messaging.RabbitMQInterface
	esClient *elastic.Client
}

func NewReportHandler(rabbitMQ messaging.RabbitMQInterface, esClient *elastic.Client) *ReportHandler {
	return &ReportHandler{rabbitMQ: rabbitMQ, esClient: esClient}
}

type ReportRequest struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

func (h *ReportHandler) RequestReport(w http.ResponseWriter, r *http.Request) {
	// Create a new report request with a unique ID and a status of "pending"
	reportRequest := ReportRequest{
		ID:     uuid.New(),
		Status: "pending",
	}

	// Publish the report request to the RabbitMQ queue
	err := h.rabbitMQ.PublishReportRequest("report_requests", reportRequest)
	if err != nil {
		// If publishing fails, return an internal server error
		http.Error(w, "Failed to request report", http.StatusInternalServerError)
		return
	}

	// Index the report request in Elasticsearch
	_, err = h.esClient.Index().
		Index("report_requests").
		Id(reportRequest.ID.String()).
		BodyJson(reportRequest).
		Do(r.Context())
	if err != nil {
		// If indexing fails, return an internal server error
		http.Error(w, "Failed to index report request", http.StatusInternalServerError)
		return
	}

	// Respond with a 202 Accepted status and encode the report request in the response
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(reportRequest)
}
