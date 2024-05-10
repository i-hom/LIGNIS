package app

import (
	"context"
	"lignis/internal/generated/api"
)

func (a App) GetCustomers(ctx context.Context, params api.GetCustomersParams) (*api.GetCustomersOK, error) {
	var response []api.CustomerWithID
	customers, total, err := a.customerRepo.GetByOption(params.Pattern.Value, int64(params.Page.Value), int64(params.Limit.Value))
	if err != nil {
		return nil, err
	}
	for i := range customers {
		response = append(response, api.CustomerWithID{
			ID:      customers[i].ID.Hex(),
			Fio:     customers[i].Fio,
			Phone:   customers[i].Phone,
			Address: customers[i].Address,
		})
	}

	return &api.GetCustomersOK{Total: int(total), Customers: response}, nil
}
