package repository

import (
	"context"
	"lignis/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SaleRepo struct {
	collection *mongo.Collection
}

func NewSaleRepo(collection *mongo.Collection) *SaleRepo {
	return &SaleRepo{
		collection: collection,
	}
}

func (s *SaleRepo) Create(sale *model.Sale) (primitive.ObjectID, error) {
	res, err := s.collection.InsertOne(context.TODO(), sale)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}
