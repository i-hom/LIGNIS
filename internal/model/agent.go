package model

import (
	"lignis/internal/generated/api"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Agent struct {
	Fio               string `bson:"fio"`
	Phone             string `bson:"phone"`
	InstagramUsername string `bson:"instagram_username,omitempty"`
	BonusPercent      uint32 `bson:"bonus_percent"`
	Is_Deleted        bool   `bson:"deleted_at,omitempty"`
}

type AgentWithID struct {
	ID    primitive.ObjectID `bson:"_id" json:"id"`
	Agent `bson:",inline"`
}

func (a AgentWithID) ToApi() api.AgentWithID {
	insta_username := api.OptString{Set: false}
	if a.InstagramUsername != "" {
		insta_username = api.OptString{Set: true, Value: a.InstagramUsername}
	}
	return api.AgentWithID{
		ID:                a.ID.Hex(),
		Fio:               a.Fio,
		Phone:             a.Phone,
		InstagramUsername: insta_username,
		BonusPercent:      int(a.BonusPercent),
	}
}
