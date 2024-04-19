package app

import (
	"context"
	"lignis/internal/generated/api"
)

func (a App) CreateCustomer(ctx context.Context, req *api.Customer) (*api.ResponseWithID, error) {
	return &api.ResponseWithID{}, nil
}
