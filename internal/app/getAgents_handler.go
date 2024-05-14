package app

import (
	"context"
	"errors"
	"lignis/internal/generated/api"
	"lignis/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a App) GetAgents(ctx context.Context, params api.GetAgentsParams) (*api.GetAgentsOK, error) {
	user := ctx.Value("user").(*model.Claims)

	if user.Role != "admin" && user.Role != "salesman" {
		return nil, errors.New("access denied")
	}

	if len(params.Pattern.Value) == 24 && params.Limit.Value == 1 {
		objectID, err := primitive.ObjectIDFromHex(params.Pattern.Value)
		if err != nil {
			return nil, err
		}
		user, err := a.agentRepo.GetByID(objectID)
		if err != nil {
			return nil, err
		}
		return &api.GetAgentsOK{Total: 1, Agents: []api.AgentWithID{{ID: user.ID.Hex(), Fio: user.Fio, Phone: user.Phone, InstagramUsername: user.InstagramUsername, BonusPercent: int(user.BonusPercent)}}}, nil
	}

	response := make([]api.AgentWithID, 0)
	agents, total, err := a.agentRepo.GetByPattern(params.Pattern.Value, int64(params.Page.Value), int64(params.Limit.Value))
	if err != nil {
		return nil, err
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
	return &api.GetAgentsOK{Total: int(total), Agents: response}, nil
}
