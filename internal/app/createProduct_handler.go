package app

import (
	"context"
	"lignis/internal/generated/api"
)

func (a App) CreateProduct(ctx context.Context, req *api.AddProductRequest) (*api.ResponseWithID, error) {
	return &api.ResponseWithID{}, nil
}
