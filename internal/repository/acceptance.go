package repository

import (
	"context"
	"lignis/internal/model"

	"go.mongodb.org/mongo-driver/bson"
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

func (a AcceptanceRepo) Delete(id primitive.ObjectID) error {
	_, err := a.collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}

func (a AcceptanceRepo) Get(id primitive.ObjectID) (*model.AcceptanceWithID, error) {
	var acceptance model.AcceptanceWithID
	err := a.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&acceptance)
	if err != nil {
		return nil, err
	}
	return &acceptance, nil
}
