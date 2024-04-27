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

func (r ProductRepo) GetByOption(pattern string, page, limit int64) ([]model.ProductWithID, error) {
	var products []model.ProductWithID

	var filter bson.M

	if len(pattern) == 24 && limit == 1 {
		id, err := primitive.ObjectIDFromHex(pattern)
		if err != nil {
			return nil, err
		}
		filter = bson.M{"_id": id}
	} else {
		filter = bson.M{
			"$or": bson.A{
				bson.M{"code": bson.M{"$regex": pattern, "$options": "i"}},
				bson.M{"name": bson.M{"$regex": pattern, "$options": "i"}},
			}, "deleted_at": bson.M{"$exists": false}}
	}

	cursor, err := r.collection.Find(
		context.TODO(),
		filter,
		options.Find().SetSort(bson.D{{Key: "name", Value: 1}}).SetLimit(limit).SetSkip((page-1)*limit))

	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (r ProductRepo) Add(product_id primitive.ObjectID, quantity uint32, price float64) error {
	err := r.collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": product_id},
		bson.M{"$inc": bson.M{"quantity": quantity},
			"$set": bson.M{"price": price}},
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
	}

	_, err = r.collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": product_id},
		bson.M{"$inc": bson.M{"quantity": -quantity}},
	)
	if err != nil {
		return err
	}
	return nil
}

func (r ProductRepo) Get(idstr string) (*model.ProductWithID, error) {
	id, err := primitive.ObjectIDFromHex(idstr)
	if err != nil {
		return nil, err
	}

	var product model.ProductWithID
	err = r.collection.FindOne(
		context.TODO(),
		bson.M{"_id": id, "deleted_at": bson.M{"$exists": false}},
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

func (r ProductRepo) Total() int64 {
	total, err := r.collection.EstimatedDocumentCount(context.TODO())
	if err != nil {
		return 0
	}
	return total
}
