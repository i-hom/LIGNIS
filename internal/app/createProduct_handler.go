package app

import (
	"context"
	"errors"
	"lignis/internal/generated/api"
	"lignis/internal/model"
)

func (a App) CreateProduct(ctx context.Context, req *api.CreateProductRequestMultipart) (*api.ResponseWithID, error) {
	user := ctx.Value("user").(*model.Claims)

	if user.Role != "manager" {
		return &api.ResponseWithID{}, errors.New("access denied")
	}

	if len(req.Name) < 5 || len(req.Code) < 5 {
		return &api.ResponseWithID{}, errors.New("invalid data")
	}

	res, err := a.productRepo.Create(&model.Product{
		Name:      req.Name,
		Code:      req.Code,
		Quantity:  0,
		SellPrice: 0,
	})

	if err != nil {
		return &api.ResponseWithID{}, err
	}
	a.minio.Upload(ctx, res.Hex(), req.Photo.File)

	return &api.ResponseWithID{ID: res.Hex()}, nil
}
