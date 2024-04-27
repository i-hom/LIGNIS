package app

import (
	"context"
	"errors"
	"lignis/internal/generated/api"
	"lignis/internal/model"
)

func (a App) UpdateProduct(ctx context.Context, req *api.EditProductRequestMultipart, params api.UpdateProductParams) (*api.SuccessResponse, error) {
	user := ctx.Value("user").(*model.Claims)

	if user.Role != "manager" {
		return &api.SuccessResponse{}, errors.New("access denied")
	}

	product, err := a.productRepo.Get(params.ID)
	if err != nil {
		return &api.SuccessResponse{}, err
	}

	if req.Name.Set {
		product.Name = req.Name.Value
	}
	if req.Code.Set {
		product.Code = req.Code.Value
	}
	if req.Price.Set {
		if product.Quantity == 0 {
			return &api.SuccessResponse{}, errors.New("there is no products in the stock")
		}
		product.SellPrice = req.Price.Value
	}
	if req.Photo.Set {
		a.minio.Delete(ctx, product.ID.Hex())
		a.minio.Upload(ctx, product.ID.Hex(), req.Photo.Value.File)
	}

	err = a.productRepo.Update(product)
	if err != nil {
		return &api.SuccessResponse{}, err
	}
	return &api.SuccessResponse{Message: "product successfully updated"}, nil
}
