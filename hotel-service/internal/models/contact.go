package models

import (
	"time"

	"github.com/google/uuid"
)

type Contact struct {
	ID        uuid.UUID `json:"id"`
	HotelID   uuid.UUID `json:"hotel_id"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
