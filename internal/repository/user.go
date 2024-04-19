package repository

import (
	"context"
	"lignis/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	collection *mongo.Collection
}

func NewUserRepo(collection *mongo.Collection) *UserRepo {
	return &UserRepo{
		collection: collection,
	}
}

func (u UserRepo) Create(user *model.User) (primitive.ObjectID, error) {
	res, err := u.collection.InsertOne(context.TODO(), user)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (u UserRepo) Get(credits model.LoginData) (*model.UserWithID, error) {
	var user model.UserWithID
	err := u.collection.FindOne(context.TODO(), bson.M{"login": credits.Login, "hashpass": credits.HashPass}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
