package app

import (
	"context"
	"errors"
	"lignis/internal/generated/api"
	"lignis/internal/model"
)

func (a App) CreateCustomer(ctx context.Context, req *api.Customer) (*api.ResponseWithID, error) {
	user := ctx.Value("user").(*model.Claims)

	if user.Role != "salesman" {
		return &api.ResponseWithID{}, errors.New("access denied")

	}

	res, err := a.customerRepo.Create(&model.Customer{
		Fio:     req.Fio,
		Phone:   req.Phone,
		Address: req.Address,
	})
	if err != nil {
		return &api.ResponseWithID{}, err
	}
	return &api.ResponseWithID{ID: res.Hex()}, nil
}
