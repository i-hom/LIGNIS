package app

import (
	"context"
	"lignis/internal/generated/api"
)

func (a App) CreateAgent(ctx context.Context, req *api.Agent) (*api.ResponseWithID, error) {
	return &api.ResponseWithID{}, nil
}
