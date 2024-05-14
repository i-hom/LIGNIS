package app

import (
	"context"
	"errors"
	"lignis/internal/generated/api"
	"lignis/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a App) DeleteUser(ctx context.Context, params api.DeleteUserParams) (*api.SuccessResponse, error) {
	user := ctx.Value("user").(*model.Claims)

	if user.Role != "admin" {
		return &api.SuccessResponse{}, errors.New("access denied")
	}

	deletionID, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		return &api.SuccessResponse{}, err
	}

	err = a.userRepo.Delete(deletionID)
	if err != nil {
		return &api.SuccessResponse{}, err
	}
	return &api.SuccessResponse{Message: "user deleted"}, nil
}
