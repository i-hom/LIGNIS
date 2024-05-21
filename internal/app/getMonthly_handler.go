package app

import (
	"context"
	"time"

	"lignis/internal/generated/api"
)

func (a App) GetMonthlyReport(ctx context.Context, params api.GetMonthlyReportParams) (*api.GetMonthlyReportOK, error) {
	date, err := time.Parse("2006-01", params.Month)

	if err != nil {
		return nil, err
	}

	data, err := a.saleRepo.GetMontylyReport(date)
	if err != nil {
		return nil, err
	}

	return &api.GetMonthlyReportOK{Report: data}, nil
}
