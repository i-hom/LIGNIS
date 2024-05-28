package app

import (
	"context"
	"errors"
	"lignis/internal/generated/api"
	"lignis/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a App) DeleteDefect(ctx context.Context, params api.DeleteDefectParams) (*api.SuccessResponse, error) {

	user := ctx.Value("user").(*model.Claims)

	if user.Role != "manager" {
		return &api.SuccessResponse{}, errors.New("access denied")
	}

	id, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		return &api.SuccessResponse{}, err
	}

	defect, err := a.defectRepo.GetByID(id)
	if err != nil {
		return &api.SuccessResponse{}, err
	}
	for _, p := range defect.Defects {
		if err != nil {
			return &api.SuccessResponse{}, err
		}
		err = a.productRepo.Add(p.ProductID, uint64(p.Quantity))
		if err != nil {
			return &api.SuccessResponse{}, err
		}
	}
	err = a.defectRepo.Delete(id)
	if err != nil {
		return &api.SuccessResponse{}, err
	}
	return &api.SuccessResponse{Message: "defect deleted"}, nil

}
