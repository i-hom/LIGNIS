package app

import (
	"context"
	"lignis/internal/generated/api"
)

func (a App) Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	if len(req.Login) < 5 || (len(req.Password) < 5 && req.Password == "") {

	}
	return &api.LoginResponse{}, nil
}
