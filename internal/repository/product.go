package repository

import (
	"context"
	"lignis/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepo struct {
	collection *mongo.Collection
}

func NewProductRepo(collection *mongo.Collection) *ProductRepo {
	return &ProductRepo{
		collection: collection,
	}
}

func (r *ProductRepo) Create(product *model.Product) (primitive.ObjectID, error) {
	res, err := r.collection.InsertOne(context.TODO(), product)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}
