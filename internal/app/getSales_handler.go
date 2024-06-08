package app

import (
	"context"
	"lignis/internal/generated/api"
)

func (a App) GetSales(ctx context.Context, params api.GetSalesParams) (*api.GetSalesOK, error) {
	isAgent := true
	isCustomer := true

	var response []api.SalesWithID
	sales, total, err := a.saleRepo.GetByDate(params.From.Value, params.To.Value, int64(params.Limit.Value), int64(params.Page.Value), false)
	if err != nil {
		return nil, err
	}

	for i := range sales {
		cart := make([]api.SaleProduct, 0)

		for j := range sales[i].Cart {
			product, err := a.productRepo.GetByID(sales[i].Cart[j].ID)
			var product_name string
			if err != nil {
				product_name = "product not found"
			}

			product_name = product.Name
			cart = append(cart, api.SaleProduct{
				ID:       sales[i].Cart[j].ID.Hex(),
				Quantity: int(sales[i].Cart[j].Quantity),
				Price:    sales[i].Cart[j].Price,
				Name:     api.NewOptString(product_name),
			})
		}
		if sales[i].AgentId.IsZero() {
			isAgent = false
		}
		if sales[i].CustomerId.IsZero() {
			isCustomer = false
		}
		response = append(response, api.SalesWithID{
			ID:           sales[i].ID.Hex(),
			AgentID:      api.OptString{Set: isAgent, Value: sales[i].AgentId.Hex()},
			CustomerID:   api.OptString{Set: isCustomer, Value: sales[i].CustomerId.Hex()},
			Date:         sales[i].Date,
			TotalUzs:     sales[i].TotalUZS,
			TotalUsd:     sales[i].TotalUSD,
			Products:     cart,
			CurrencyCode: api.SalesWithIDCurrencyCode(sales[i].CurrencyCode),
		})
	}
	return &api.GetSalesOK{Total: int(total), Sales: response}, nil
}
