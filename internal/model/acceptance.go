package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Acceptance struct {
	Products   []ShortProduct     `bson:"products"`
	AcceptedBy primitive.ObjectID `bson:"accepted_by"`
}

type AcceptanceWithID struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	Acceptance `bson:",inline"`
}
