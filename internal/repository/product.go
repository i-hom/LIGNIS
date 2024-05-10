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

type ProductRepo struct {
	collection *mongo.Collection
}

func NewProductRepo(collection *mongo.Collection) *ProductRepo {
	return &ProductRepo{
		collection: collection,
	}
}

func (r ProductRepo) Create(product *model.Product) (primitive.ObjectID, error) {
	var dbproduct model.ProductWithID
	err := r.collection.FindOne(context.TODO(), bson.M{"code": product.Code}).Decode(&dbproduct)
	if err == nil {
		return primitive.NilObjectID, errors.New("product already exists")
	}

	res, err := r.collection.InsertOne(context.TODO(), product)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (r ProductRepo) GetByOption(pattern string, page, limit int64) ([]model.ProductWithID, int64, error) {
	var products []model.ProductWithID

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
				bson.M{"code": bson.M{"$regex": pattern, "$options": "i"}},
				bson.M{"name": bson.M{"$regex": pattern, "$options": "i"}},
			}, "is_deleted": bson.M{"$exists": false}}
	}

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

	if err := cursor.All(context.TODO(), &products); err != nil {
		return nil, 0, err
	}
	return products, count, nil
}

func (r ProductRepo) Add(product_id primitive.ObjectID, quantity uint32, price float64) error {
	err := r.collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": product_id},
		bson.M{"$inc": bson.M{"quantity": quantity},
			"$set": bson.M{"sell_price": price}},
	).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r ProductRepo) Consume(product_id primitive.ObjectID, quantity uint32) error {
	var product model.ProductWithID
	err := r.collection.FindOne(
		context.TODO(),
		bson.M{"_id": product_id},
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

func (r ProductRepo) Get(id primitive.ObjectID) (*model.ProductWithID, error) {
	var product model.ProductWithID
	err := r.collection.FindOne(
		context.TODO(),
		bson.M{"_id": id, "id_deleted": bson.M{"$exists": false}},
	).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r ProductRepo) Update(product *model.ProductWithID) error {
	_, err := r.collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": product.ID},
		bson.M{"$set": product},
	)
	if err != nil {
		return err
	}
	return nil
}

func (r ProductRepo) GetStats() (int, int, float64, error) {
	var total_quantity int
	var total_stock_value float64
	var result bson.M

	cursor, err := r.collection.Aggregate(
		context.TODO(),
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "total_quantity", Value: bson.D{
				{Key: "$sum", Value: "$quantity"},
			},
			},
		},
		},
		})

	if err != nil {
		return 0, 0, 0, err
	}
	defer cursor.Close(context.TODO())

	cursor.All(context.TODO(), &result)
	total_quantity = result["total_quantity"].(int)

	cursor, err = r.collection.Aggregate(
		context.TODO(),
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "totalCostQuantity", Value: bson.D{
					{Key: "$sum", Value: bson.D{
						{Key: "$multiply", Value: bson.A{"$cost", "$quantity"}},
					}},
				}},
			}},
		})
	if err != nil {
		return 0, 0, 0, err
	}

	cursor.All(context.TODO(), &result)
	total_stock_value = result["totalCostQuantity"].(float64)

	total_products, err := r.collection.EstimatedDocumentCount(context.TODO())
	if err != nil {
		return 0, 0, 0, err
	}

	return int(total_products), int(total_quantity), total_stock_value, nil
}
