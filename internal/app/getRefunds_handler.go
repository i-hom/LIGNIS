package app

import (
	"context"
	"errors"
	"lignis/internal/generated/api"
	"lignis/internal/model"
)

type ByObject struct {
	Name    string
	IsExist bool
}

func (a App) GetRefunds(ctx context.Context, params api.GetRefundsParams) (*api.GetRefundsOK, error) {
	user := ctx.Value("user").(*model.Claims)
	agent := ByObject{"deleted agent", false}
	customer := ByObject{"deleted customer", false}
	deleter := ByObject{"deleted user", false}
	response := make([]api.SalesWithID, 0)
	if user.Role != "manager" && user.Role != "admin" {
		return nil, errors.New("access denied")
	}

	refund, total, err := a.saleRepo.GetByDate(params.From.Value, params.To.Value, int64(params.Limit.Value), int64(params.Page.Value), true)
	if err != nil {
		return nil, err
	}
	for i := range refund {
		cart := make([]api.SaleProduct, 0)

		for j := range refund[i].Cart {
			product, err := a.productRepo.GetByID(refund[i].Cart[j].ID)
			var product_name string
			if err != nil {
				product_name = "product not found"
			}

			product_name = product.Name
			cart = append(cart, api.SaleProduct{
				ID:       refund[i].Cart[j].ID.Hex(),
				Quantity: int(refund[i].Cart[j].Quantity),
				Price:    refund[i].Cart[j].Price,
				Name:     api.NewOptString(product_name),
			})
		}
		if !refund[i].AgentId.IsZero() {
			agent.IsExist = true
			AgentName, err := a.agentRepo.GetByID(refund[i].AgentId)
			if err == nil {
				agent.Name = AgentName.Fio
			}
		}
		if !refund[i].CustomerId.IsZero() {
			customer.IsExist = true
			CustomerName, err := a.customerRepo.GetByID(refund[i].CustomerId)
			if err == nil {
				customer.Name = CustomerName.Fio
			}
		}
		if !refund[i].Deleted_By.IsZero() {
			deleter.IsExist = true
			DeletedBy, err := a.userRepo.GetByID(refund[i].Deleted_By)
			if err == nil {
				deleter.Name = DeletedBy.Fio
			}
		}
		response = append(response, api.SalesWithID{
			ID:           refund[i].ID.Hex(),
			AgentID:      api.OptString{Set: agent.IsExist, Value: agent.Name},
			CustomerID:   api.OptString{Set: customer.IsExist, Value: customer.Name},
			DeletedBy:    api.OptString{Set: deleter.IsExist, Value: deleter.Name},
			Date:         refund[i].Date,
			TotalUzs:     refund[i].TotalUZS,
			TotalUsd:     refund[i].TotalUSD,
			Products:     cart,
			CurrencyCode: api.SalesWithIDCurrencyCode(refund[i].CurrencyCode),
		})
	}

	return &api.GetRefundsOK{Total: int(total), Refunds: response}, nil
}
