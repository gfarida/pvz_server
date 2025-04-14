package model

import "time"

type ProductType string

const (
	Electronics ProductType = "электроника"
	Clothing    ProductType = "одежда"
	Shoes       ProductType = "обувь"
)

var AllowedProductTypes = map[ProductType]bool{
	"электроника": true,
	"одежда":      true,
	"обувь":       true,
}

type Product struct {
	ID          string      `json:"id"`
	DateTime    time.Time   `json:"dateTime"`
	Type        ProductType `json:"type"`
	ReceptionID string      `json:"receptionId"`
}
