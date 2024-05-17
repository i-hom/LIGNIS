package app

import (
	"context"
	"errors"
	"lignis/internal/generated/api"
	"lignis/internal/model"
)

func (a App) GetArchivedProducts(ctx context.Context, params api.GetArchivedProductsParams) (*api.GetArchivedProductsOK, error) {
	user := ctx.Value("user").(*model.Claims)
	var response []api.GetProducts

	if user.Role != "manager" {
		return nil, errors.New("access denied")
	}
	products, total, err := a.productRepo.GetByPattern(params.Pattern.Value, int64(params.Page.Value), int64(params.Limit.Value), true)
	if err != nil {
		return nil, err
	}

	for i := range products {
		response = append(response, api.GetProducts{
			ID:       products[i].ID.Hex(),
			Name:     products[i].Name,
			Code:     products[i].Code,
			Quantity: int(products[i].Quantity),
			Price:    products[i].SellPrice,
		})
	}
	return &api.GetArchivedProductsOK{Total: int(total), Products: response}, nil
}
