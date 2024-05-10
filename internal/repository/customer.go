package repository

import (
	"context"
	"errors"
	"lignis/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CustomerRepo struct {
	collection *mongo.Collection
}

func NewCustomerRepo(collection *mongo.Collection) *CustomerRepo {
	return &CustomerRepo{
		collection: collection,
	}
}

func (r CustomerRepo) Create(customer *model.Customer) (primitive.ObjectID, error) {
	var dbcustomer model.CustomerWithID
	err := r.collection.FindOne(context.TODO(), bson.M{"phone": customer.Phone}).Decode(&dbcustomer)
	if err == nil {
		return primitive.NilObjectID, errors.New("customer already exists")
	}
	res, err := r.collection.InsertOne(context.TODO(), customer)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (r CustomerRepo) GetByOption(pattern string, page, limit int64) ([]model.CustomerWithID, int64, error) {
	var customers []model.CustomerWithID

	var filter bson.M

	if len(pattern) == 24 && limit == 1 {
		id, err := primitive.ObjectIDFromHex(pattern)
		if err != nil {
			return nil, 0, err
		}
		filter = bson.M{"_id": id}
	} else {
		filter = bson.M{
			"$or": bson.A{
				bson.M{"address": bson.M{"$regex": pattern, "$options": "i"}},
				bson.M{"phone": bson.M{"$regex": pattern, "$options": "i"}},
				bson.M{"fio": bson.M{"$regex": pattern, "$options": "i"}},
			}}
	}

	cursor, err := r.collection.Find(
		context.TODO(),
		filter,
		options.Find().SetSort(bson.D{{Key: "fio", Value: 1}}).SetLimit(limit).SetSkip((page-1)*limit))
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.TODO())

	count, err := r.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, 0, err
	}

	if err := cursor.All(context.TODO(), &customers); err != nil {
		return nil, 0, err
	}
	return customers, count, nil
}
