package entity

import (
	"time"
)

type PVZ struct {
	ID               string    `json:"id" db:"id"`
	RegistrationDate time.Time `json:"registration_date" db:"registration_date"`
	City             string    `json:"city" db:"city"`
}
