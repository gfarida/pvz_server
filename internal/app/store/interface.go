package store

import (
	"context"
	"pvz_server/internal/app/model"
	"time"
)

type PVZCreator interface {
	CreatePVZ(ctx context.Context, city model.City) (*model.PVZ, error)
}

type ReceptionCreator interface {
	CreateReception(ctx context.Context, pvzID string) (*model.Reception, error)
}

type ProductAdder interface {
	AddProduct(ctx context.Context, pvzID string, prodType model.ProductType) (*model.Product, error)
}

type ProductDeleter interface {
	DeleteLastProduct(ctx context.Context, pvzID string) error
}

type ReceptionCloser interface {
	CloseLastReception(ctx context.Context, pvzID string) (*model.Reception, error)
}

type PVZFetcher interface {
	FetchPVZList(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]*model.PVZWithReceptions, error)
}
