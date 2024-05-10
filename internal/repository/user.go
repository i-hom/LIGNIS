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
	var dbuser model.UserWithID
	err := u.collection.FindOne(
		context.TODO(), bson.M{
			"$or": bson.A{
				bson.M{"login": user.Login},
				bson.M{"fio": user.Fio},
			},
		}).Decode(&dbuser)
	if err == nil {
		return primitive.NilObjectID, errors.New("user already exists")
	}
	res, err := u.collection.InsertOne(context.TODO(), user)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (u UserRepo) GetByOption(pattern string, page, limit int64) ([]model.UserWithID, int64, error) {
	var users []model.UserWithID
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
				bson.M{"login": bson.M{"$regex": pattern, "$options": "i"}},
				bson.M{"fio": bson.M{"$regex": pattern, "$options": "i"}},
			}, "deleted_at": bson.M{"$ne": nil}}
	}
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
	if err = cursor.All(context.TODO(), &users); err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (u UserRepo) Get(credits model.LoginData) (*model.UserWithID, error) {
	var user model.UserWithID
	err := u.collection.FindOne(context.TODO(), bson.M{"login": credits.Login, "hashpass": credits.HashPass}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
