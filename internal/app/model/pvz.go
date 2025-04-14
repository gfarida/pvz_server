package model

import (
	"time"
)

type City string

const (
	Moscow City = "Москва"
	SPB    City = "Санкт-Петербург"
	Kazan  City = "Казань"
)

var AllowedCities = map[City]bool{
	Moscow: true,
	SPB:    true,
	Kazan:  true,
}

type PVZ struct {
	ID               string    `json:"id"`
	RegistrationDate time.Time `json:"registrationDate"`
	City             City      `json:"city"`
}

type PVZWithReceptions struct {
	PVZ        PVZ                     `json:"pvz"`
	Receptions []ReceptionWithProducts `json:"receptions"`
}

type ReceptionWithProducts struct {
	Reception Reception `json:"reception"`
	Products  []Product `json:"products"`
}
