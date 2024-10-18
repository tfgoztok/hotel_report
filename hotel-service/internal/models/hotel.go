package models

import (
	"time"

	"github.com/google/uuid"
)

type Hotel struct {
	ID              uuid.UUID `json:"id"`
	OfficialName    string    `json:"official_name"`
	OfficialSurname string    `json:"official_surname"`
	CompanyTitle    string    `json:"company_title"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
