package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Customer struct {
	Name    string `bson:"name"`
	Phone   string `bson:"phone"`
	Address string `bson:"address"`
}

type CustomerWithID struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Customer `bson:",inline"`
}
