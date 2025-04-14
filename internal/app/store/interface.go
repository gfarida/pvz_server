package store

import (
	"context"
	"pvz_server/internal/app/model"
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
