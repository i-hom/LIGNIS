package repository

import (
	"context"
	"fmt"
	"lignis/internal/generated/api"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MonthlyRepo struct {
	collection *mongo.Collection
}

func NewMonthlyRepo(collection *mongo.Collection) *MonthlyRepo {
	return &MonthlyRepo{
		collection: collection,
	}
}

func (r MonthlyRepo) Operation(date time.Time, value float64) error {

	year, month, _ := date.Date()
	err := r.collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{
			"month": month,
			"year":  year,
		},
		bson.M{
			"$inc": bson.M{"value": value},
		},
	)
	if err.Err() != nil {
		_, err := r.collection.InsertOne(context.TODO(), bson.M{
			"month": month,
			"year":  year,
			"value": value,
		})
		return err
	}
	return nil
}

func (r MonthlyRepo) GetYearly() ([]api.Analytics, error) {
	var analytics []api.Analytics

	year, _, _ := time.Now().Date()

	cursor, err := r.collection.Find(context.TODO(), bson.M{"year": year}, options.Find().SetSort(bson.D{{Key: "_id", Value: 1}}))
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		var a bson.M
		err := cursor.Decode(&a)
		if err != nil {
			return nil, err
		}
		analytics = append(analytics, api.Analytics{
			Label: fmt.Sprintf("%d-%d", a["month"].(int32), a["year"].(int32)),
			Value: a["value"].(float64),
		})
	}
	return analytics, nil
}
