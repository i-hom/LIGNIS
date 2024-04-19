package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type LoginData struct {
	Login    string `bson:"login"`
	HashPass string `bson:"hashpass"`
}

type User struct {
	Name      string `bson:"name"`
	LoginData `bson:",inline"`
	Role      string `bson:"role"`
}

type UserWithID struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	User `bson:",inline"`
}
