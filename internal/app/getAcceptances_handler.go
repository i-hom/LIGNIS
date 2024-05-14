package app

import (
	"context"
	"errors"
	"lignis/internal/generated/api"
	"lignis/internal/model"
)

func (a App) GetAcceptances(ctx context.Context, params api.GetAcceptancesParams) (*api.GetAcceptancesOK, error) {
	user := ctx.Value("user").(*model.Claims)

	response := make([]api.AcceptanceWithID, 0)

	if user.Role != "manager" {
		return nil, errors.New("access denied")
	}

	acceptances, count, err := a.acceptanceRepo.GetByDate(params.From.Value, params.To.Value, int64(params.Limit.Value), int64(params.Page.Value))
	if err != nil {
		return nil, err
	}

	for i := range acceptances {

		products := make([]api.AcceptanceProduct, 0)

		for p := range acceptances[i].Products {
			product, _ := a.productRepo.GetByID(acceptances[i].Products[p].ID)
			products = append(products, api.AcceptanceProduct{
				ID:        acceptances[i].Products[p].ID.Hex(),
				Name:      product.Name,
				Quantity:  int(acceptances[i].Products[p].Quantity),
				CostPrice: acceptances[i].Products[p].Price,
			})
		}

		response = append(response, api.AcceptanceWithID{
			ID:       acceptances[i].ID.Hex(),
			Date:     acceptances[i].Date,
			Products: products,
		})
	}
	return &api.GetAcceptancesOK{Total: int(count), Acceptances: response}, nil
}
