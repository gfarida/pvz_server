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
	ErrProductTypeNotAllowed  = errors.New("unsupported product type")
	ErrNoActiveReception      = errors.New("no active reception for this PVZ")
)

func (s *Store) CreatePVZ(ctx context.Context, city model.City) (*model.PVZ, error) {
	if !model.AllowedCities[city] {
		return nil, ErrCityNotAllowed
	}

	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, ErrDatabase
	}

	defer tx.Rollback()

	id := uuid.NewString()
	now := time.Now()

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO pvz (id, registration_date, city) 
		VALUES ($1, $2, $3)`,
		id,
		now,
		city,
	)

	if err != nil {
		return nil, ErrDatabase
	}

	if err := tx.Commit(); err != nil {
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
					SELECT 1 FROM reception
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

	if err := tx.Commit(); err != nil {
		return nil, ErrDatabase
	}

	return &model.Reception{
		ID:       id,
		DateTime: now,
		PvzID:    pvzID,
		Status:   model.InProgress,
	}, nil
}

func (s *Store) AddProduct(ctx context.Context, pvzID string, productType model.ProductType) (*model.Product, error) {
	if !model.AllowedProductTypes[productType] {
		return nil, ErrProductTypeNotAllowed
	}

	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, ErrDatabase
	}

	defer tx.Rollback()

	var receptionID string
	err = tx.QueryRowContext(
		ctx,
		`SELECT id FROM reception
		 WHERE pvz_id = $1 AND status = $2
		 ORDER BY date_time DESC
		 LIMIT 1`,
		pvzID,
		model.InProgress,
	).Scan(&receptionID)

	if err != nil {
		return nil, ErrNoActiveReception
	}

	id := uuid.NewString()
	now := time.Now()

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO product (id, date_time, type, reception_id)
		 VALUES ($1, $2, $3, $4)`,
		id,
		now,
		productType,
		receptionID,
	)

	if err != nil {
		return nil, ErrDatabase
	}

	if err := tx.Commit(); err != nil {
		return nil, ErrDatabase
	}

	return &model.Product{
		ID:          id,
		DateTime:    now,
		Type:        productType,
		ReceptionID: receptionID,
	}, nil
}
