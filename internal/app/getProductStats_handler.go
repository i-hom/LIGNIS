package app

import (
	"context"
	"lignis/internal/generated/api"
)

func (a App) GetProductStats(ctx context.Context) (*api.GetProductStatsOK, error) {
	tp, tq, tsv, err := a.productRepo.GetStats()
	return &api.GetProductStatsOK{TotalQuantity: tq, TotalProducts: tp, TotalStockValue: tsv}, err
}
