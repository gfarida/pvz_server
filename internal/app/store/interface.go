package store

import (
	"context"
	"pvz_server/internal/app/model"
)

type PVZCreator interface {
	CreatePVZ(ctx context.Context, city model.City) (*model.PVZ, error)
}
