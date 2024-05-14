package app

import (
	"context"
	"errors"
	"lignis/internal/generated/api"
	"lignis/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a App) DeleteAgent(ctx context.Context, params api.DeleteAgentParams) (*api.SuccessResponse, error) {

	user := ctx.Value("user").(*model.Claims)

	if user.Role != "admin" {
		return &api.SuccessResponse{}, errors.New("access denied")
	}

	deletionID, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		return &api.SuccessResponse{}, err
	}

	err = a.agentRepo.Delete(deletionID)
	if err != nil {
		return &api.SuccessResponse{}, err
	}
	return &api.SuccessResponse{Message: "agent deleted"}, nil
}
