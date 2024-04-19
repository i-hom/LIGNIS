package repository

import (
	"context"
	"lignis/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AcceptanceRepo struct {
	collection *mongo.Collection
}

func NewAcceptanceRepo(collection *mongo.Collection) *AcceptanceRepo {
	return &AcceptanceRepo{
		collection: collection,
	}
}

func (a AcceptanceRepo) Create(acceptance *model.Acceptance) (primitive.ObjectID, error) {
	res, err := a.collection.InsertOne(context.TODO(), acceptance)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}
