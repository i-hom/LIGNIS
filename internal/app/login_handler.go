package app

import (
	"context"
	"errors"
	"lignis/internal/generated/api"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func (a App) Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	if len(req.Login) < 5 || len(req.Password) < 5 || strings.Contains(req.Password, " ") {
		return &api.LoginResponse{}, errors.New("invalid data")
	}

	user, err := a.userRepo.GetByLogin(req.Login)
	if err != nil {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashPass), []byte(req.Password))
	if err != nil {
		return nil, err
	}

	token, err := a.auth.GenerateToken(user)
	if err != nil {
		return &api.LoginResponse{}, errors.New("login failed")
	}

	return &api.LoginResponse{Token: token}, nil
}
