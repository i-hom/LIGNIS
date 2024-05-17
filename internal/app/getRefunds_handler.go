package app

import (
	"context"
	"errors"
	"lignis/internal/generated/api"
	"lignis/internal/model"
)

func (a App) GetRefunds(ctx context.Context, params api.GetRefundsParams) (*api.GetRefundsOK, error) {
	user := ctx.Value("user").(*model.Claims)

	if user.Role != "manager" && user.Role != "admin" {
		return nil, errors.New("access denied")
	}

	refund, total, err := a.saleRepo.GetByDate(params.From.Value, params.To.Value, int64(params.Limit.Value), int64(params.Page.Value), true)
	if err != nil {
		return nil, err
	}
	response := make([]api.SalesWithID, 0)
	for i := range refund {

		response = append(response, api.SalesWithID{
			ID: refund[i].ID.Hex(),
		})
	}
	return &api.GetRefundsOK{Total: int(total), Refunds: response}, nil
}
