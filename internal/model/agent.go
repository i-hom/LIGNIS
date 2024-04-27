package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Agent struct {
	Fio               string             `bson:"fio"`
	Phone             string             `bson:"phone"`
	InstagramUsername string             `bson:"instagram_username"`
	BonusPercent      uint32             `bson:"bonus_percent"`
	Deleted_At        primitive.DateTime `bson:"deleted_at,omitempty"`
}

type AgentWithID struct {
	ID    primitive.ObjectID `bson:"_id" json:"id"`
	Agent `bson:",inline"`
}
