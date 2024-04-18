package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Name     string `bson:"name"`
	Login    string `bson:"login"`
	HashPass string `bson:"hashpass"`
	Role     string `bson:"role"`
}

type UserWithID struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	User `bson:",inline"`
}
