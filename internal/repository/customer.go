package repository

import (
	"context"
	"lignis/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerRepo struct {
	collection *mongo.Collection
}

func NewCustomerRepo(collection *mongo.Collection) *CustomerRepo {
	return &CustomerRepo{
		collection: collection,
	}
}

func (r *CustomerRepo) Create(customer *model.Customer) (primitive.ObjectID, error) {
	res, err := r.collection.InsertOne(context.TODO(), customer)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}
