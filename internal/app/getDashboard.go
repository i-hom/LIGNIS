package app

import (
	"context"
	"lignis/internal/generated/api"
)

func (a App) GetDashboard(ctx context.Context) (*api.Dashboard, error) {

	topProduct, err := a.saleRepo.GetThisMonthTopSoldProduct(10)
	if err != nil {
		return nil, err
	}

	weekly, err := a.saleRepo.GetLastWeek()
	if err != nil {
		return nil, err
	}

	yearly, err := a.monthlyRepo.GetYearly()
	if err != nil {
		return nil, err
	}

	return &api.Dashboard{
		Last7daySales: weekly,
		TopProducts:   topProduct,
		LastYearSales: yearly,
	}, nil
}
