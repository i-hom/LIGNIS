package app

import (
	"context"
	"lignis/internal/generated/api"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a App) GetCustomers(ctx context.Context, params api.GetCustomersParams) (*api.GetCustomersOK, error) {
	if len(params.Pattern.Value) == 24 && params.Limit.Value == 1 {
		objectID, err := primitive.ObjectIDFromHex(params.Pattern.Value)
		if err != nil {
			return nil, err
		}
		customer, err := a.customerRepo.GetByID(objectID)
		if err != nil {
			return nil, err
		}
		return &api.GetCustomersOK{Total: 1, Customers: []api.CustomerWithID{{ID: customer.ID.Hex(), Fio: customer.Fio, Phone: customer.Phone, Address: customer.Address}}}, nil
	}

	response := make([]api.CustomerWithID, 0)
	customers, total, err := a.customerRepo.GetByPattern(params.Pattern.Value, int64(params.Page.Value), int64(params.Limit.Value))
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
