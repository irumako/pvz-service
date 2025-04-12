package entity

import (
	"time"
)

type ReceptionStatus string

const (
	InProgress ReceptionStatus = "in_progress"
	Close      ReceptionStatus = "close"
)

type Reception struct {
	ID       string          `json:"id" db:"id"`
	Datetime time.Time       `json:"dateTime" db:"datetime"`
	PvzID    string          `json:"pvzId" db:"pvz_id"`
	Status   ReceptionStatus `json:"status" db:"status"`
}
