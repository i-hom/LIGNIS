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

type UserRepo struct {
	collection *mongo.Collection
}

func NewUserRepo(collection *mongo.Collection) *UserRepo {
	return &UserRepo{
		collection: collection,
	}
}

func (u UserRepo) Create(user *model.User) (primitive.ObjectID, error) {
	if u.collection.FindOne(context.TODO(), bson.M{"login": user.Login}).Err() == nil {
		return primitive.NilObjectID, errors.New("user already exists")
	}

	res, err := u.collection.InsertOne(context.TODO(), user)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (u UserRepo) GetByPattern(pattern string, page, limit int64) ([]model.UserWithID, int64, error) {
	var users []model.UserWithID
	filter := bson.M{
		"$or": bson.A{
			bson.M{"login": bson.M{"$regex": pattern, "$options": "i"}},
			bson.M{"fio": bson.M{"$regex": pattern, "$options": "i"}},
			bson.M{"role": bson.M{"$regex": pattern, "$options": "i"}},
		}, "is_deleted": bson.M{"$exists": false}}

	cursor, err := u.collection.Find(
		context.TODO(),
		filter,
		options.Find().SetSort(bson.D{{Key: "fio", Value: 1}}).SetLimit(limit).SetSkip((page-1)*limit))

	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.TODO())

	count, err := u.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, 0, err
	}

	for cursor.Next(context.TODO()) {
		var user model.UserWithID
		err = cursor.Decode(&user)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}
	return users, count, nil
}

func (u UserRepo) GetByID(id primitive.ObjectID) (*model.UserWithID, error) {
	var user model.UserWithID
	err := u.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u UserRepo) GetByLogin(login string) (*model.UserWithID, error) {
	var user model.UserWithID
	err := u.collection.FindOne(context.TODO(), bson.M{"login": login, "is_deleted": bson.M{"$exists": false}}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u UserRepo) Delete(id primitive.ObjectID) error {
	err := u.collection.FindOneAndUpdate(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"is_deleted": true}})
	return err.Err()
}
