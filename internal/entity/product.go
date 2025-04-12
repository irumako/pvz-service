package entity

import (
	"time"
)

type ProductType string

const (
	Electronic ProductType = "электроника"
	Clothing   ProductType = "одежда"
	Footwear   ProductType = "обувь"
)

type Product struct {
	ID          string      `json:"id" db:"id"`
	Datetime    time.Time   `json:"dateTime" db:"datetime"`
	ReceptionID string      `json:"receptionId" db:"reception_id"`
	Type        ProductType `json:"type" db:"type"`
}
