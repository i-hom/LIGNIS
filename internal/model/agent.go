package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Agent struct {
	Name              string  `bson:"name"`
	Phone             string  `bson:"phone"`
	InstagramUsername string  `bson:"instagram_username"`
	BonusPercent      float64 `bson:"bonus_percent"`
}

type AgentWithID struct {
	ID    primitive.ObjectID `bson:"_id" json:"id"`
	Agent `bson:",inline"`
}
