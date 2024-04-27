package app

import (
	"context"
	"errors"
	"lignis/internal/generated/api"
	"lignis/internal/model"
)

func (a App) GetAgents(ctx context.Context, params api.GetAgentsParams) ([]api.AgentWithID, error) {
	user := ctx.Value("user").(*model.Claims)

	if user.Role != "admin" {
		return []api.AgentWithID{}, errors.New("access denied")
	}

	var response []api.AgentWithID
	agents, err := a.agentRepo.GetByOption(params.Pattern.Value, int64(params.Page.Value), int64(params.Limit.Value))
	if err != nil {
		return []api.AgentWithID{}, err
	}

	for i := range agents {
		response = append(response, api.AgentWithID{
			ID:                agents[i].ID.Hex(),
			Fio:               agents[i].Fio,
			Phone:             agents[i].Phone,
			InstagramUsername: agents[i].InstagramUsername,
			BonusPercent:      int(agents[i].BonusPercent),
		})
	}
	return response, nil
}
