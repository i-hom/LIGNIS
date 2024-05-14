package app

import (
	"context"
	"errors"
	"lignis/internal/generated/api"
	"lignis/internal/model"
)

func (a App) CreateAgent(ctx context.Context, req *api.Agent) (*api.ResponseWithID, error) {
	user := ctx.Value("user").(*model.Claims)

	if user.Role != "admin" {
		return &api.ResponseWithID{}, errors.New("access denied")
	}

	res, err := a.agentRepo.Create(&model.Agent{
		Fio:               req.Fio,
		Phone:             req.Phone,
		InstagramUsername: req.InstagramUsername,
		BonusPercent:      uint32(req.BonusPercent),
	})
	if err != nil {
		return &api.ResponseWithID{}, err
	}
	return &api.ResponseWithID{ID: res.Hex()}, nil
}
