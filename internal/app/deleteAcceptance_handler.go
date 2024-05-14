package app

import (
	"context"
	"errors"
	"lignis/internal/generated/api"
	"lignis/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a App) DeleteAcceptance(ctx context.Context, params api.DeleteAcceptanceParams) (*api.SuccessResponse, error) {
	user := ctx.Value("user").(*model.Claims)

	if user.Role != "manager" {
		return &api.SuccessResponse{}, errors.New("access denied")
	}

	id, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		return &api.SuccessResponse{}, err
	}

	acceptance, err := a.acceptanceRepo.GetByID(id)
	if err != nil {
		return &api.SuccessResponse{}, err
	}
	for _, p := range acceptance.Products {
		err = a.productRepo.Consume(p.ID, p.Quantity)
		if err != nil {
			return &api.SuccessResponse{}, err
		}
	}

	err = a.acceptanceRepo.Delete(id)
	if err != nil {
		return &api.SuccessResponse{}, err
	}
	return &api.SuccessResponse{Message: "acceptance deleted"}, nil
}
