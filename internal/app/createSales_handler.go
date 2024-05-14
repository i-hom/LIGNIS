package app

import (
	"context"
	"errors"
	"lignis/internal/generated/api"
	"lignis/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a App) CreateSale(ctx context.Context, req *api.Sales) (*api.ResponseWithID, error) {
	user := ctx.Value("user").(*model.Claims)

	if user.Role != "salesman" {
		return &api.ResponseWithID{}, errors.New("access denied")
	}

	if len(req.Products) < 1 {
		return &api.ResponseWithID{}, errors.New("there is no products in the cart")
	}

	var sales []model.ShortProduct
	var newSale model.Sale
	var err error

	for _, p := range req.Products {
		if p.Quantity < 1 {
			return &api.ResponseWithID{}, errors.New("invalid quantity")
		}

		id, err := primitive.ObjectIDFromHex(p.ID)
		if err != nil {
			return &api.ResponseWithID{}, err
		}

		err = a.productRepo.Consume(id, uint32(p.Quantity))
		if err != nil {
			return &api.ResponseWithID{}, err
		}

		sales = append(sales, model.ShortProduct{
			ID:       id,
			Quantity: uint32(p.Quantity),
			Price:    p.Price,
		})
	}

	newSale = model.Sale{
		SalesmanId:   user.UserID,
		Cart:         sales,
		TotalUZS:     req.TotalUzs,
		TotalUSD:     req.TotalUsd,
		CurrencyCode: string(req.CurrencyCode),
	}

	if req.CustomerID.Set {
		newSale.CustomerId, err = primitive.ObjectIDFromHex(req.CustomerID.Value)
		if err != nil {
			return &api.ResponseWithID{}, errors.New("invalid customer id")
		}
	}

	if req.AgentID.Set {
		newSale.AgentId, err = primitive.ObjectIDFromHex(req.AgentID.Value)
		if err != nil {
			return &api.ResponseWithID{}, errors.New("invalid agent id")
		}
	}

	res, err := a.saleRepo.Create(&newSale)
	if err != nil {
		return &api.ResponseWithID{}, err
	}
	return &api.ResponseWithID{ID: res.Hex()}, nil
}
