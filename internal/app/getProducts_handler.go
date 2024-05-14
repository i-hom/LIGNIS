package app

import (
	"context"
	"lignis/internal/generated/api"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a App) GetProducts(ctx context.Context, params api.GetProductsParams) (*api.GetProductsOK, error) {

	if len(params.Pattern.Value) == 24 && params.Limit.Value == 1 {
		objectID, err := primitive.ObjectIDFromHex(params.Pattern.Value)
		if err != nil {
			return nil, err
		}
		product, err := a.productRepo.GetByID(objectID)
		if err != nil {
			return nil, err
		}
		return &api.GetProductsOK{Total: 1, Products: []api.GetProducts{{ID: product.ID.Hex(), Name: product.Name, Code: product.Code, Quantity: int(product.Quantity), Price: product.SellPrice}}}, nil
	}

	products, total, err := a.productRepo.GetByPattern(params.Pattern.Value, int64(params.Page.Value), int64(params.Limit.Value))
	response := make([]api.GetProducts, 0)
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

	return &api.GetProductsOK{Total: int(total), Products: response}, nil
}
