package model

import "time"

type ProductType string

const (
	Electronics ProductType = "электроника"
	Clothing    ProductType = "одежда"
	Shoes       ProductType = "обувь"
)

type Product struct {
	ID          string      `json:"id"`
	DateTime    time.Time   `json:"dateTime"`
	Type        ProductType `json:"type"`
	ReceptionID string      `json:"receptionId"`
}
