package store

import (
	"context"
	"errors"
	"pvz_server/internal/app/model"
	"time"

	"github.com/google/uuid"
)

var (
	ErrCityNotAllowed = errors.New("unsupported city")
	ErrDatabase       = errors.New("database error")
)

func (s *Store) CreatePVZ(ctx context.Context, city model.City) (*model.PVZ, error) {
	if !model.AllowedCities[city] {
		return nil, ErrCityNotAllowed
	}

	id := uuid.New()
	now := time.Now()

	_, err := s.db.ExecContext(
		ctx,
		`INSERT INTO pvz (id, registration_date, city) VALUES ($1, $2, $3)`,
		id,
		now,
		city,
	)

	if err != nil {
		return nil, ErrDatabase
	}

	return &model.PVZ{
		ID:               id,
		RegistrationDate: now,
		City:             city,
	}, nil
}
