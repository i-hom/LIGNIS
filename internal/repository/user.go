package repository

import (
	"context"
	"lignis/internal/model"

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

func (u *UserRepo) Create(user *model.User) (primitive.ObjectID, error) {
	res, err := u.collection.InsertOne(context.TODO(), user)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}
