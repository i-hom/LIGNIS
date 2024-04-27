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

type AgentRepo struct {
	collection *mongo.Collection
}

func NewAgentRepo(collection *mongo.Collection) *AgentRepo {
	return &AgentRepo{
		collection: collection,
	}
}

func (a AgentRepo) Create(agent *model.Agent) (primitive.ObjectID, error) {
	var dbagent model.AgentWithID
	err := a.collection.FindOne(
		context.TODO(), bson.M{
			"$or": bson.A{
				bson.M{"phone": agent.Phone},
				bson.M{"instagram_username": agent.InstagramUsername},
			},
		}).Decode(&dbagent)
	if err == nil {
		return primitive.NilObjectID, errors.New("agent already exists")
	}
	res, err := a.collection.InsertOne(context.TODO(), agent)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (a AgentRepo) GetByOption(pattern string, page, limit int64) ([]model.AgentWithID, error) {
	var agents []model.AgentWithID
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
				bson.M{"phone": bson.M{"$regex": pattern, "$options": "i"}},
				bson.M{"fio": bson.M{"$regex": pattern, "$options": "i"}},
				bson.M{"instagram_username": bson.M{"$regex": pattern, "$options": "i"}},
			}, "deleted_at": bson.M{"$ne": nil}}
	}

	cursor, err := a.collection.Find(
		context.TODO(),
		filter,
		options.Find().SetSort(bson.D{{Key: "fio", Value: 1}}).SetLimit(limit).SetSkip((page-1)*limit))

	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &agents); err != nil {
		return nil, err
	}
	return agents, nil
}

func (a AgentRepo) Total() int64 {
	total, err := a.collection.EstimatedDocumentCount(context.TODO())
	if err != nil {
		return 0
	}
	return total
}
