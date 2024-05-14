package app

import (
	"context"
	"errors"
	"lignis/internal/generated/api"
	"lignis/internal/model"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func (a App) CreateUser(ctx context.Context, req *api.CreateUserRequest) (*api.ResponseWithID, error) {
	user := ctx.Value("user").(*model.Claims)

	if user.Role != "admin" {
		return &api.ResponseWithID{}, errors.New("access denied")
	}

	if strings.Contains(req.Password, " ") || strings.Contains(req.Login, " ") {
		return &api.ResponseWithID{}, errors.New("invalid data")
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return &api.ResponseWithID{}, err
	}

	res, err := a.userRepo.Create(&model.User{
		Fio: req.Fio,
		LoginData: model.LoginData{
			Login:    req.Login,
			HashPass: string(pass),
		},
		Role: string(req.Role),
	})
	if err != nil {
		return &api.ResponseWithID{}, err
	}

	return &api.ResponseWithID{ID: res.Hex()}, nil
}
