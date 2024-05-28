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
		{{Key: "$match", Value: bson.M{"is_deleted": bson.M{"$exists": false}, "_id": bson.M{"$gte": primitive.NewObjectIDFromTimestamp(thisMonth), "$lt": primitive.NewObjectIDFromTimestamp(nextMonth)}}}},
		{{Key: "$unwind", Value: "$cart"}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$cart._id"},
			{Key: "value", Value: bson.D{{Key: "$sum", Value: "$cart.quantity"}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "value", Value: -1}}}},
		{{Key: "$limit", Value: 10}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "products"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "product"},
		}}},
		{{Key: "$unwind", Value: "$product"}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "label", Value: "$product.name"},
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
	var monthly_usd float64 = 0.0
	var monthly_uzs int64 = 0
	cursor, err := s.collection.Aggregate(
		context.TODO(),
		[]bson.M{
			{
				"$match": bson.M{"is_deleted": bson.M{"$exists": false}, "_id": bson.M{"$gte": primitive.NewObjectIDFromTimestamp(month), "$lt": primitive.NewObjectIDFromTimestamp(month.AddDate(0, 1, 0))}},
			},
			{
				"$group": bson.M{
					"_id": bson.M{
						"date": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d %H:%M", "date": bson.M{
							"$dateAdd": bson.M{
								"startDate": bson.M{
									"$toDate": "$_id",
								},
								"unit":   "hour",
								"amount": 5,
							},
						}}},
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
			monthly_uzs += result["total_uzs"].(int64)
		}

		if currency == "USD" {
			data[date].TotalUsd = result["total_usd"].(float64)
			monthly_usd += result["total_usd"].(float64)
		}
	}

	for _, v := range data {
		stats = append(stats, *v)
	}

	stats = append(stats, api.DailyReport{
		Date:     "Total",
		TotalUzs: monthly_uzs,
		TotalUsd: monthly_usd,
	})
	return stats, nil
}

func (s SaleRepo) GetMonthlyBonus(month time.Time, id primitive.ObjectID) (api.GetAgentReportOK, error) {
	stats := make([]api.GetAgentReportOKReportItem, 0)

	cursor, err := s.collection.Aggregate(context.TODO(), []bson.M{
		{
			"$match": bson.M{
				"_id": bson.M{
					"$gte": primitive.NewObjectIDFromTimestamp(month),
					"$lt":  primitive.NewObjectIDFromTimestamp(month.AddDate(0, 1, 0)),
				},
				"agent_id": id,
			},
		},
		{
			"$group": bson.M{
				"_id":   "$agent_id",
				"total": bson.M{"$sum": "$total_usd"},
				"recipts": bson.M{
					"$push": bson.M{
						"id": "$_id",
						"date": bson.M{
							"$dateToString": bson.M{
								"format": "%Y-%m-%d %H:%M", "date": bson.M{
									"$dateAdd": bson.M{
										"startDate": bson.M{
											"$toDate": "$_id",
										},
										"unit":   "hour",
										"amount": 5,
									},
								},
							},
						},
						"total": "$total_usd",
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "agents",
				"localField":   "_id",
				"foreignField": "_id",
				"as":           "agent",
			},
		},
		{"$unwind": "$agent"},
		{
			"$addFields": bson.M{
				"bonus": bson.M{"$multiply": []interface{}{"$total", bson.M{"$multiply": []interface{}{0.01, "$agent.bonus_percent"}}}},
			},
		},
		{
			"$project": bson.M{
				"bonus":   1,
				"recipts": 1,
				"total":   1,
			},
		},
	})

	if err != nil {
		return api.GetAgentReportOK{}, err
	}
	defer cursor.Close(context.TODO())

	result := make([]bson.M, 0)
	err = cursor.All(context.TODO(), &result)

	if err != nil {
		return api.GetAgentReportOK{}, err
	}

	if len(result) == 0 {
		return api.GetAgentReportOK{}, nil
	}
	recipts := result[0]["recipts"].(primitive.A)
	for _, v := range recipts {
		stats = append(stats, api.GetAgentReportOKReportItem{
			ID:       v.(bson.M)["id"].(primitive.ObjectID).Hex(),
			Date:     v.(bson.M)["date"].(string),
			TotalUsd: v.(bson.M)["total"].(float64),
		})
	}

	return api.GetAgentReportOK{
		Report: stats,
		Total:  result[0]["total"].(float64),
		Bonus:  result[0]["bonus"].(float64),
	}, nil
}
