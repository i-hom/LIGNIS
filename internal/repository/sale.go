package repository

import (
	"context"
	"lignis/internal/generated/api"
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

func (s SaleRepo) GetByDate(from, to string, limit, page int64, is_deleted bool) ([]model.SaleWithID, int64, error) {
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
			}, "is_deleted": bson.M{"$exists": is_deleted},
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
		sale.Date = sale.ID.Timestamp().Add(time.Hour * 5)
		sales = append(sales, sale)
	}
	return sales, count, nil
}

func (s SaleRepo) GetThisMonthTopSoldProduct(limit int64) ([]api.Analytics, error) {
	var products []api.Analytics
	thisMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
	nextMonth := thisMonth.AddDate(0, 1, 0)
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"is_deleted": false, "_id": bson.M{"$gte": primitive.NewObjectIDFromTimestamp(thisMonth), "$lt": primitive.NewObjectIDFromTimestamp(nextMonth)}}}},
		{{Key: "$unwind", Value: "$cart"}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$cart._id"},
			{Key: "value", Value: bson.D{{Key: "$sum", Value: "$cart.quantity"}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "value", Value: -1}}}},
		{{Key: "$limit", Value: limit}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "products"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "product"},
		}}},
		{{Key: "$unwind", Value: "$product"}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "key", Value: "$product.name"},
			{Key: "value", Value: 1},
		}}},
	}

	cursor, err := s.collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		var product api.Analytics
		err := cursor.Decode(&product)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (s SaleRepo) Delete(id, deleted_by primitive.ObjectID) error {
	result := s.collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"is_deleted": true, "deleted_by": deleted_by}},
	)
	return result.Err()
}

func (s SaleRepo) GetLastWeek() ([]api.Analytics, error) {
	stats := make([]api.Analytics, 0)
	lastWeek := time.Now().AddDate(0, 0, -7)
	cursor, err := s.collection.Aggregate(
		context.TODO(),
		[]bson.M{
			{
				"$match": bson.M{"is_deleted": bson.M{"$exists": false}, "_id": bson.M{"$gte": primitive.NewObjectIDFromTimestamp(lastWeek), "$lte": primitive.NewObjectIDFromTimestamp(time.Now())}},
			},
			{
				"$group": bson.M{
					"_id":       bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": bson.M{"$toDate": "$_id"}}},
					"total_usd": bson.M{"$sum": "$total_usd"},
				},
			},
		})

	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		result := bson.M{}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}

		stats = append(stats, api.Analytics{
			Label: result["_id"].(string),
			Value: result["total_usd"].(float64),
		})
	}

	return stats, nil
}

func (s SaleRepo) GetMontylyReport(month time.Time) ([]api.DailyReport, error) {
	stats := make([]api.DailyReport, 0)
	cursor, err := s.collection.Aggregate(
		context.TODO(),
		[]bson.M{
			{
				"$match": bson.M{"is_deleted": bson.M{"$exists": false}, "_id": bson.M{"$gte": primitive.NewObjectIDFromTimestamp(month), "$lt": primitive.NewObjectIDFromTimestamp(month.AddDate(0, 1, 0))}},
			},
			{
				"$group": bson.M{
					"_id": bson.M{
						"date":          bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": bson.M{"$toDate": "$_id"}}},
						"currency_code": "$currency_code",
					},
					"total_uzs": bson.M{"$sum": "$total_uzs"},
					"total_usd": bson.M{"$sum": "$total_usd"},
				},
			},
		})

	if err != nil {
		return nil, err
	}
	data := make(map[string]*api.DailyReport, 0)
	for cursor.Next(context.TODO()) {
		result := bson.M{}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}

		_id := result["_id"].(bson.M)
		date := _id["date"].(string)
		currency := _id["currency_code"].(string)

		if data[date] == nil {
			data[date] = &api.DailyReport{
				Date:     date,
				TotalUzs: 0,
				TotalUsd: 0,
			}
		}

		if currency == "UZS" {
			data[date].TotalUzs = result["total_uzs"].(int64)
		}

		if currency == "USD" {
			data[date].TotalUsd = result["total_usd"].(float64)
		}
	}

	for _, v := range data {
		stats = append(stats, *v)
	}
	return stats, nil
}
