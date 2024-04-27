package app

import (
	"context"
	"lignis/internal/generated/api"
)

func (a App) GetProducts(ctx context.Context, params api.GetProductsParams) ([]api.GetProducts, error) {
	products, err := a.productRepo.GetByOption(params.Pattern.Value, int64(params.Page.Value), int64(params.Limit.Value))
	response := make([]api.GetProducts, 0)
	if err != nil {
		return []api.GetProducts{}, err
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

	return response, nil
}
