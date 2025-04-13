package store

import (
	"context"
	"errors"
	"pvz_server/internal/app/model"
	"time"

	"github.com/google/uuid"
)

var (
	ErrCityNotAllowed         = errors.New("unsupported city")
	ErrDatabase               = errors.New("database error")
	ErrReceptionAlreadyExists = errors.New("receprion in progress")
)

func (s *Store) CreatePVZ(ctx context.Context, city model.City) (*model.PVZ, error) {
	if !model.AllowedCities[city] {
		return nil, ErrCityNotAllowed
	}

	id := uuid.NewString()
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

func (s *Store) CreateReception(ctx context.Context, pvzID string) (*model.Reception, error) {
	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, ErrDatabase
	}

	defer tx.Rollback()

	var exists bool

	err = tx.QueryRowContext(
		ctx,
		`SELECT EXISTS (
					SELECT 1 FROM receptionct
					WHERE pvz_id = $1 AND status = $2
				)`,
		pvzID,
		model.InProgress,
	).Scan(
		&exists,
	)

	if err != nil {
		return nil, ErrDatabase
	}

	if exists {
		return nil, ErrReceptionAlreadyExists
	}

	id := uuid.NewString()
	now := time.Now()

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO reception (id, date_time, pvz_id, status)
		VALUES ($1, $2, $3, $4)`,
		id,
		now,
		pvzID,
		model.InProgress,
	)

	if err != nil {
		return nil, ErrDatabase
	}

	return &model.Reception{
		ID:       id,
		DateTime: now,
		PvzID:    pvzID,
		Status:   model.InProgress,
	}, nil
}
