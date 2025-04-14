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
	ErrNoProductsToDelete     = errors.New("no products to delete")
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

func (s *Store) DeleteLastProduct(ctx context.Context, pvzID string) error {
	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return ErrDatabase
	}

	defer tx.Rollback()

	var receptionID string
	err = tx.QueryRowContext(ctx,
		`SELECT id FROM reception
		WHERE pvz_id = $1 AND status = $2
		ORDER BY date_time DESC LIMIT 1`,
		pvzID,
		model.InProgress,
	).Scan(&receptionID)

	if err != nil {
		return ErrNoActiveReception
	}

	var productID string

	err = tx.QueryRowContext(ctx,
		`SELECT id FROM product
		WHERE reception_id = $1
		ORDER BY date_time DESC LIMIT 1`,
		receptionID,
	).Scan(&productID)

	if err != nil {
		return ErrNoProductsToDelete
	}

	_, err = tx.ExecContext(ctx,
		`DELETE FROM product 
		WHERE id = $1`,
		productID,
	)

	if err != nil {
		return ErrDatabase
	}

	return tx.Commit()
}

func (s *Store) CloseLastReception(ctx context.Context, pvzID string) (*model.Reception, error) {
	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, ErrDatabase
	}

	defer tx.Rollback()

	var r model.Reception

	err = tx.QueryRowContext(
		ctx,
		`SELECT id, date_time, status FROM reception
		 WHERE pvz_id = $1 AND status = $2
		 ORDER BY date_time DESC LIMIT 1`,
		pvzID,
		model.InProgress,
	).Scan(
		&r.ID,
		&r.DateTime,
		&r.Status,
	)

	if err != nil {
		return nil, ErrNoActiveReception
	}

	_, err = tx.ExecContext(
		ctx,
		`UPDATE reception SET status = $1 
		WHERE id = $2`,
		model.Closed,
		r.ID,
	)

	if err != nil {
		return nil, ErrDatabase
	}

	if err := tx.Commit(); err != nil {
		return nil, ErrDatabase
	}

	r.PvzID = pvzID
	r.Status = model.Closed
	return &r, nil
}

func (s *Store) FetchPVZList(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]*model.PVZWithReceptions, error) {
	offset := (page - 1) * limit

	query :=
		`SELECT p.id, p.registration_date, p.city, 
		       r.id, r.date_time, r.status,
		       pr.id, pr.date_time, pr.type
		FROM pvz p
		LEFT JOIN reception r ON r.pvz_id = p.id
		LEFT JOIN product pr ON pr.reception_id = r.id
		WHERE ($1::timestamp IS NULL OR r.date_time >= $1)
		  AND ($2::timestamp IS NULL OR r.date_time <= $2)
		ORDER BY p.registration_date
		OFFSET $3 LIMIT $4`

	rows, err := s.db.QueryContext(
		ctx,
		query,
		startDate,
		endDate,
		offset,
		limit,
	)

	if err != nil {
		return nil, ErrDatabase
	}

	defer rows.Close()

	result, err := aggregatePVZResults(rows)

	if err != nil {
		return nil, ErrDatabase
	}

	return result, nil
}
