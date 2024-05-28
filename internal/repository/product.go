package repository

import (
	"context"
	"errors"
	"lignis/internal/model"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepo struct {
	collection *mongo.Collection
}

func NewProductRepo(collection *mongo.Collection) *ProductRepo {
	return &ProductRepo{
		collection: collection,
	}
}

func (r ProductRepo) Create(product *model.Product) (primitive.ObjectID, error) {
	var p model.ProductWithID

	err := r.collection.FindOne(
		context.TODO(),
		bson.M{"code": product.Code},
	).Decode(&p)
	if err == nil {
		if p.Is_Deleted {
			res, err := r.collection.UpdateOne(context.TODO(), bson.M{"_id": p.ID}, product)
			return res.UpsertedID.(primitive.ObjectID), err
		}
		if !p.Is_Deleted {
			return primitive.NilObjectID, errors.New("product already exists")
		}
	}
	res, err := r.collection.InsertOne(context.TODO(), product)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (r ProductRepo) Update(product *model.ProductWithID) error {
	_, err := r.collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": product.ID, "is_deleted": bson.M{"$exists": false}},
		bson.M{"$set": product},
	)
	if err != nil {
		return err
	}
	return nil
}

func (r ProductRepo) GetByPattern(pattern string, page, limit int64, is_deleted bool) ([]model.ProductWithID, int64, error) {
	var products []model.ProductWithID

	filter := bson.M{
		"$or": bson.A{
			bson.M{"code": bson.M{"$regex": pattern, "$options": "i"}},
			bson.M{"name": bson.M{"$regex": pattern, "$options": "i"}},
		}, "is_deleted": bson.M{"$exists": is_deleted}}

	cursor, err := r.collection.Find(
		context.TODO(),
		filter,
		options.Find().SetSort(bson.D{{Key: "name", Value: 1}}).SetLimit(limit).SetSkip((page-1)*limit))

	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.TODO())

	count, err := r.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, 0, err
	}

	for cursor.Next(context.TODO()) {
		var product model.ProductWithID
		err := cursor.Decode(&product)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, product)
	}
	return products, count, nil
}

func (r ProductRepo) Add(product_id primitive.ObjectID, quantity uint64) error {
	err := r.collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": product_id, "is_deleted": bson.M{"$exists": false}},
		bson.M{"$inc": bson.M{"quantity": quantity}},
	).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r ProductRepo) Consume(product_id primitive.ObjectID, quantity uint64) error {
	var product model.ProductWithID
	err := r.collection.FindOne(
		context.TODO(),
		bson.M{"_id": product_id, "is_deleted": bson.M{"$exists": false}},
	).Decode(&product)

	if err != nil {
		return err
	}

	if product.Quantity < quantity {
		return errors.New("not enough products")
	} else {
		product.Quantity -= quantity
		err = r.collection.FindOneAndUpdate(
			context.TODO(),
			bson.M{"_id": product_id},
			bson.M{"$set": product},
		).Err()
	}

	if err != nil {
		return err
	}
	return nil
}

func (r ProductRepo) GetStats() (int, int, float64, error) {
	var total_quantity int
	var total_stock_value float64
	var result []bson.M

	filter := bson.M{"is_deleted": bson.M{"$exists": false}}

	cursor, err := r.collection.Aggregate(context.TODO(),
		[]bson.M{
			{"$match": filter},
			{"$group": bson.M{
				"_id": nil,
				"tq":  bson.M{"$sum": "$quantity"},
				"tsv": bson.M{"$sum": bson.M{"$multiply": bson.A{"$sell_price", "$quantity"}}}},
			}},
	)

	if err != nil {
		return 0, 0, 0, err
	}
	err = cursor.All(context.TODO(), &result)
	if err != nil {
		return 0, 0, 0, err
	}

	if len(result) == 0 {
		return 0, 0, 0, nil
	}
	if result[0]["tq"] == nil || result[0]["tsv"] == nil {
		return 0, 0, 0, nil
	}

	if reflect.TypeOf(result[0]["tq"]) == reflect.TypeOf(int32(0)) {
		total_quantity = int(result[0]["tq"].(int32))
	} else if reflect.TypeOf(result[0]["tq"]) == reflect.TypeOf(int64(0)) {
		total_quantity = int(result[0]["tq"].(int64))
	}
	total_stock_value, ok := result[0]["tsv"].(float64)
	if !ok {
		total_stock_value = 0
	}
	total_products, err := r.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, 0, 0, err
	}
	return int(total_products), total_quantity, total_stock_value, nil
}

func (r ProductRepo) GetByID(id primitive.ObjectID) (*model.ProductWithID, error) {
	var product model.ProductWithID
	err := r.collection.FindOne(
		context.TODO(),
		bson.M{"_id": id},
	).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r ProductRepo) Delete(id primitive.ObjectID) error {
	err := r.collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": id, "quantity": bson.M{"$eq": 0}},
		bson.M{"$set": bson.M{"is_deleted": true}},
	).Err()
	return err
}

func (r ProductRepo) Restore(product_id primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(context.TODO(), bson.M{"_id": product_id, "is_deleted": true}, bson.M{"$unset": bson.M{"is_deleted": 1}})
	if err != nil {
		return err
	}
	return nil
}
