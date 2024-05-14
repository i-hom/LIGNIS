package app

import (
	"context"
	"errors"
	"lignis/internal/generated/api"
	"lignis/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a App) GetUsers(ctx context.Context, params api.GetUsersParams) (*api.GetUsersOK, error) {
	user := ctx.Value("user").(*model.Claims)

	response := make([]api.UserWithID, 0)

	if user.Role != "admin" {
		return nil, errors.New("access denied")
	}

	if len(params.Pattern.Value) == 24 && params.Limit.Value == 1 {
		objectID, err := primitive.ObjectIDFromHex(params.Pattern.Value)
		if err != nil {
			return nil, err
		}
		user, err := a.userRepo.GetByID(objectID)
		if err != nil {
			return nil, err
		}
		return &api.GetUsersOK{Total: 1, Users: []api.UserWithID{{ID: user.ID.Hex(), Fio: user.Fio, Role: api.UserWithIDRole(user.Role)}}}, nil
	} else {

		users, total, err := a.userRepo.GetByPattern(params.Pattern.Value, int64(params.Page.Value), int64(params.Limit.Value))
		if err != nil {
			return nil, err
		}

		for i := range users {
			response = append(response, api.UserWithID{
				ID:   users[i].ID.Hex(),
				Fio:  users[i].Fio,
				Role: api.UserWithIDRole(users[i].Role),
			})
		}
		return &api.GetUsersOK{Total: int(total), Users: response}, nil
	}
}
