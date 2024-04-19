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

func (a App) CreateUser(ctx context.Context, req *api.User) (*api.ResponseWithID, error) {
	user := ctx.Value("user").(*model.Claims)

	if user.Role != "admin" {
		return &api.ResponseWithID{}, errors.New("access denied")
	}

	if len(req.Fio) < 5 || len(req.Login) < 5 || len(req.Password) < 5 || strings.Contains(req.Password, " ") {
		return &api.ResponseWithID{}, errors.New("invalid data")
	}

	hasher := fnv.New32a()
	_, err := hasher.Write([]byte(req.Password))
	if err != nil {
		return &api.ResponseWithID{}, err
	}
	pass := hasher.Sum32()

	res, err := a.userRepo.Create(&model.User{
		Name: req.Fio,
		LoginData: model.LoginData{
			Login:    req.Login,
			HashPass: fmt.Sprint(pass),
		},
		Role: req.Role,
	})
	if err != nil {
		return &api.ResponseWithID{}, err
	}
	return &api.ResponseWithID{ID: res.Hex()}, nil
}
