package app

import (
	"context"
	"errors"
	"lignis/internal/generated/api"
	"lignis/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a App) GetAgentReport(ctx context.Context, params api.GetAgentReportParams) (*api.GetAgentReportOK, error) {

	user := ctx.Value("user").(*model.Claims)

	if user.Role != "admin" && user.Role != "salesman" {
		return nil, errors.New("access denied")
	}

	date, err := time.Parse("2006-01", params.Month)

	if err != nil {
		return nil, err
	}

	agentID, err := primitive.ObjectIDFromHex(params.ID)

	if err != nil {
		return nil, err
	}

	data, err := a.saleRepo.GetMonthlyBonus(date, agentID)

	if err != nil {
		return nil, err
	}

	return &data, nil
}
