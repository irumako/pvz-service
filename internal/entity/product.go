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

func ToDBProductType(pt ProductType) string {
	switch pt {
	case Electronic:
		return "electronics"
	case Clothing:
		return "clothes"
	case Footwear:
		return "shoes"
	default:
		return ""
	}
}

func FromDBProductType(s string) ProductType {
	switch s {
	case "electronics":
		return Electronic
	case "clothes":
		return Clothing
	case "shoes":
		return Footwear
	default:
		return ""
	}
}

type Product struct {
	ID          string      `json:"id" db:"id"`
	Datetime    time.Time   `json:"dateTime" db:"datetime"`
	ReceptionID string      `json:"receptionId" db:"reception_id"`
	Type        ProductType `json:"type" db:"type"`
}
