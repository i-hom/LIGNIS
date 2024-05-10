package app

import (
	"context"
	"errors"
	"lignis/internal/generated/api"
	"lignis/internal/model"
)

func (a App) GetSales(ctx context.Context, params api.GetSalesParams) (*api.GetSalesOK, error) {
	user := ctx.Value("user").(*model.Claims)

	if user.Role != "admin" && user.Role != "salesman" {
		return nil, errors.New("access denied")
	}

	var response []api.SalesWithID
	sales, total, err := a.saleRepo.GetByDate(params.From.Value, params.To.Value, int64(params.Limit.Value), int64(params.Page.Value))
	if err != nil {
		return nil, err
	}

	for i := range sales {
		cart := make([]api.SaleProduct, 0)

		for j := range sales[i].Cart {
			product, err := a.productRepo.Get(sales[i].Cart[j].ID)
			var product_name string
			if err != nil {
				product_name = "product not found"
			}

			product_name = product.Name
			cart = append(cart, api.SaleProduct{
				ID:       sales[i].Cart[j].ID.Hex(),
				Quantity: int(sales[i].Cart[j].Quantity),
				Price:    sales[i].Cart[j].Price,
				Name:     api.OptString{Set: true, Value: product_name},
			})
		}

		response = append(response, api.SalesWithID{
			ID:           sales[i].ID.Hex(),
			AgentID:      sales[i].AgentId.Hex(),
			CustomerID:   sales[i].CustomerId.Hex(),
			Date:         sales[i].Date,
			TotalUzs:     sales[i].TotalUZS,
			TotalUsd:     sales[i].TotalUSD,
			Products:     cart,
			CurrencyCode: api.SalesWithIDCurrencyCode(sales[i].CurrencyCode),
		})
	}
	return &api.GetSalesOK{Total: int(total), Sales: response}, nil
}
