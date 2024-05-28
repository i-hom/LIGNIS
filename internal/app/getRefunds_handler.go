package app

import (
	"context"
	"errors"
	"lignis/internal/generated/api"
	"lignis/internal/model"
)

func (a App) GetRefunds(ctx context.Context, params api.GetRefundsParams) (*api.GetRefundsOK, error) {
	user := ctx.Value("user").(*model.Claims)
	isAgent := true
	isCustomer := true
	isDeletedBy := true
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
		if refund[i].AgentId.IsZero() {
			isAgent = false
		}
		if refund[i].CustomerId.IsZero() {
			isCustomer = false
		}
		if refund[i].Deleted_By.IsZero() {
			isDeletedBy = false
		}
		response = append(response, api.SalesWithID{
			ID:           refund[i].ID.Hex(),
			AgentID:      api.OptString{Set: isAgent, Value: refund[i].AgentId.Hex()},
			CustomerID:   api.OptString{Set: isCustomer, Value: refund[i].CustomerId.Hex()},
			DeletedBy:    api.OptString{Set: isDeletedBy, Value: refund[i].Deleted_By.Hex()},
			Date:         refund[i].Date,
			TotalUzs:     refund[i].TotalUZS,
			TotalUsd:     refund[i].TotalUSD,
			Products:     cart,
			CurrencyCode: api.SalesWithIDCurrencyCode(refund[i].CurrencyCode),
		})
	}

	return &api.GetRefundsOK{Total: int(total), Refunds: response}, nil
}
