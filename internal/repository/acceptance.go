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

func (a AcceptanceRepo) GetByDate(from, to string, limit, page int64) ([]model.AcceptanceWithID, int64, error) {
	var acceptances []model.AcceptanceWithID
	var filter bson.M
	format := "2006-01-02"

	if from != "" && to != "" {
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
			},
		}
	} else {
		filter = bson.M{}
	}
	cursor, err := a.collection.Find(
		context.TODO(),
		filter,
		options.Find().SetSort(bson.D{{Key: "_id", Value: 1}}).SetLimit(limit).SetSkip((page-1)*limit))
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.TODO())

	count, err := a.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, 0, err
	}

	for cursor.Next(context.TODO()) {
		var result model.AcceptanceWithID
		err := cursor.Decode(&result)
		if err != nil {
			return nil, 0, err
		}
		result.Date = result.ID.Timestamp()
		acceptances = append(acceptances, result)
	}
	return acceptances, count, err
}

func (a AcceptanceRepo) GetByID(id primitive.ObjectID) (*model.AcceptanceWithID, error) {
	var acceptance model.AcceptanceWithID
	err := a.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&acceptance)
	return &acceptance, err
}
