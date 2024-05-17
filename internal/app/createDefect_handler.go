package app

import (
	"context"
	"errors"
	"lignis/internal/generated/api"
	"lignis/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a App) CreateDefect(ctx context.Context, req *api.CreateDefectReq) (*api.ResponseWithID, error) {
	user := ctx.Value("user").(*model.Claims)

	if user.Role != "salesman" && user.Role != "manager" {
		return &api.ResponseWithID{}, errors.New("access denied")
	}

	defects := make([]model.DefectProduct, 0)

	for _, p := range req.Products {
		id, err := primitive.ObjectIDFromHex(p.ProductID)
		if err != nil {
			return &api.ResponseWithID{}, err
		}
		err = a.productRepo.Consume(id, uint32(p.Quantity))
		if err != nil {
			return &api.ResponseWithID{}, err
		}

		defects = append(defects, model.DefectProduct{
			ProductID: id,
			Quantity:  uint32(p.Quantity),
			Remark:    p.Remark.Value,
		})
	}

	res, err := a.defectRepo.Create(&model.Defect{
		CreatedBy: user.UserID,
		Defects:   defects,
	})
	if err != nil {
		return &api.ResponseWithID{}, err
	}
	return &api.ResponseWithID{ID: res.Hex()}, nil
}
