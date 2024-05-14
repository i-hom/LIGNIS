package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type LoginData struct {
	Login    string `bson:"login"`
	HashPass string `bson:"hashpass"`
}

type User struct {
	Fio        string `bson:"fio"`
	LoginData  `bson:",inline"`
	Role       string `bson:"role"`
	Is_Deleted bool   `bson:"deleted_at,omitempty"`
}

type UserWithID struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	User `bson:",inline"`
}
