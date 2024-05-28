package app

import (
	"context"
	"errors"
	"lignis/internal/generated/api"
	"lignis/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a App) DeleteArchivedProducts(ctx context.Context, params api.DeleteArchivedProductsParams) (*api.SuccessResponse, error) {
	user := ctx.Value("user").(*model.Claims)

	if user.Role != "manager" {
		return &api.SuccessResponse{}, errors.New("access denied")
	}

	productID, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		return &api.SuccessResponse{}, err
	}

	err = a.productRepo.Restore(productID)
	if err != nil {
		return &api.SuccessResponse{}, err
	}
	return &api.SuccessResponse{Message: "product restored"}, nil
}
