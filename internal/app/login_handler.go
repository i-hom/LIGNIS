package app

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"lignis/internal/generated/api"
	"lignis/internal/model"
	"strings"
)

func (a App) Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	if len(req.Login) < 5 || len(req.Password) < 5 || strings.Contains(req.Password, " ") {
		return &api.LoginResponse{}, errors.New("invalid data")
	}

	hasher := fnv.New32a()
	_, err := hasher.Write([]byte(req.Password))
	if err != nil {
		return &api.LoginResponse{}, err
	}
	pass := hasher.Sum32()

	user, err := a.userRepo.Get(model.LoginData{Login: req.Login, HashPass: fmt.Sprint(pass)})
	if err != nil {
		return &api.LoginResponse{}, err
	}

	token, err := a.auth.GenerateToken(user)
	if err != nil {
		return &api.LoginResponse{}, err
	}

	return &api.LoginResponse{Token: token}, nil
}
