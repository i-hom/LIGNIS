package app

import (
	"context"
	"errors"

	"lignis/internal/generated/api"
	"lignis/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a App) DeleteSale(ctx context.Context, params api.DeleteSaleParams) (*api.SuccessResponse, error) {
	user := ctx.Value("user").(*model.Claims)

	if user.Role != "salesman" {
		return &api.SuccessResponse{}, errors.New("access denied")
	}

	deletionID, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		return &api.SuccessResponse{}, err
	}

	sale, err := a.saleRepo.GetByID(deletionID)
	if err != nil {
		return &api.SuccessResponse{}, err
	}

	for _, p := range sale.Cart {
		err = a.productRepo.Add(p.ID, p.Quantity)
		if err != nil {
			return &api.SuccessResponse{}, err
		}
	}

	err = a.saleRepo.Delete(deletionID, user.UserID)
	if err != nil {
		return &api.SuccessResponse{}, err
	}

	return &api.SuccessResponse{Message: "sale deleted"}, nil
}
