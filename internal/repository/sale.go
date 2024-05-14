package repository

import (
	"context"
	"lignis/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SaleRepo struct {
	collection *mongo.Collection
}

func NewSaleRepo(collection *mongo.Collection) *SaleRepo {
	return &SaleRepo{
		collection: collection,
	}
}

func (s SaleRepo) Create(sale *model.Sale) (primitive.ObjectID, error) {
	res, err := s.collection.InsertOne(context.TODO(), sale)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (s SaleRepo) GetByID(id primitive.ObjectID) (*model.SaleWithID, error) {
	var sale model.SaleWithID
	err := s.collection.FindOne(
		context.TODO(),
		bson.M{"_id": id},
	).Decode(&sale)
	return &sale, err
}

func (s SaleRepo) GetByDate(from, to string, limit, page int64) ([]model.SaleWithID, int64, error) {
	var sales []model.SaleWithID
	var filter bson.M
	format := "2006-01-02"

	if from == "" || to == "" {
		filter = bson.M{}
	} else {

		tfrom, err := time.Parse(format, from)
		if err != nil {
			return nil, 0, err
		}
		tto, err := time.Parse(format, to)
		if err != nil {
			return nil, 0, err
		}

		if tfrom == tto {
			tto = tto.Add(24 * time.Hour)
		}

		filter = bson.M{
			"_id": bson.M{
				"$gte": primitive.NewObjectIDFromTimestamp(tfrom),
				"$lte": primitive.NewObjectIDFromTimestamp(tto),
			}, "is_deleted": bson.M{"$exists": false},
		}
	}
	cursor, err := s.collection.Find(
		context.TODO(),
		filter,
		options.Find().SetSort(bson.D{{Key: "_id", Value: 1}}).SetLimit(limit).SetSkip((page-1)*limit))
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.TODO())

	count, err := s.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, 0, err
	}

	for cursor.Next(context.TODO()) {
		var sale model.SaleWithID
		err := cursor.Decode(&sale)
		if err != nil {
			return nil, 0, err
		}
		sale.Date = sale.ID.Timestamp()
		sales = append(sales, sale)
	}
	return sales, count, nil
}

func (s SaleRepo) Delete(id primitive.ObjectID) error {
	result := s.collection.FindOneAndUpdate(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"is_deleted": true}})
	return result.Err()
}
