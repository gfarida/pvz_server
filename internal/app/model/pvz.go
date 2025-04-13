package model

import (
	"time"

	"github.com/google/uuid"
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
	ID               uuid.UUID `json:"id"`
	RegistrationDate time.Time `json:"registrationDate"`
	City             City      `json:"city"`
}
