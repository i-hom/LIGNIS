package app

import (
	"context"
	"errors"
	"lignis/internal/generated/api"
	"lignis/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a App) CreateAcceptance(ctx context.Context, req *api.CreateAcceptanceReq) (*api.ResponseWithID, error) {
	user := ctx.Value("user").(*model.Claims)

	if user.Role != "manager" {
		return &api.ResponseWithID{}, errors.New("access denied")
	}

	if len(req.Products) < 1 {
		return &api.ResponseWithID{}, errors.New("invalid data")
	}

	var acceptance []model.ShortProduct

	for _, p := range req.Products {
		id, err := primitive.ObjectIDFromHex(p.ID)
		if err != nil {
			return &api.ResponseWithID{}, err
		}

		product, err := a.productRepo.GetByID(id)
		if err != nil {
			return &api.ResponseWithID{}, err
		}

		product.Quantity += uint64(p.Quantity)
		product.SellPrice = p.SellPrice

		err = a.productRepo.Update(product)
		if err != nil {
			return &api.ResponseWithID{}, err
		}

		acceptance = append(acceptance, model.ShortProduct{
			ID:       id,
			Quantity: uint64(p.Quantity),
			Price:    p.CostPrice,
		})
	}

	res, err := a.acceptanceRepo.Create(&model.Acceptance{
		AcceptedBy: user.UserID,
		Products:   acceptance,
	})
	if err != nil {
		return &api.ResponseWithID{}, err
	}

	return &api.ResponseWithID{ID: res.Hex()}, nil
}
