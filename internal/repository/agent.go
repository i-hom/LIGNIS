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
	if a.collection.FindOne(
		context.TODO(), bson.M{
			"$or": bson.A{
				bson.M{"phone": agent.Phone},
				bson.M{"instagram_username": agent.InstagramUsername},
			},
		}).Err() == nil {
		return primitive.NilObjectID, errors.New("agent already exists")
	}
	res, err := a.collection.InsertOne(context.TODO(), agent)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (a AgentRepo) GetByPattern(pattern string, page, limit int64) ([]model.AgentWithID, int64, error) {
	var agents []model.AgentWithID
	filter := bson.M{
		"$or": bson.A{
			bson.M{"phone": bson.M{"$regex": pattern, "$options": "i"}},
			bson.M{"fio": bson.M{"$regex": pattern, "$options": "i"}},
			bson.M{"instagram_username": bson.M{"$regex": pattern, "$options": "i"}},
		}, "is_deleted": bson.M{"$exists": false}}

	cursor, err := a.collection.Find(
		context.TODO(),
		filter,
		options.Find().SetSort(bson.D{{Key: "fio", Value: 1}}).SetLimit(limit).SetSkip((page-1)*limit))

	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.TODO())

	count, err := a.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, 0, err
	}

	for cursor.Next(context.TODO()) {
		var agent model.AgentWithID
		err = cursor.Decode(&agent)
		if err != nil {
			return nil, 0, err
		}
		agents = append(agents, agent)
	}

	return agents, count, nil
}

func (a AgentRepo) GetByID(id primitive.ObjectID) (*model.AgentWithID, error) {
	var agent model.AgentWithID
	err := a.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&agent)
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func (a AgentRepo) Delete(id primitive.ObjectID) error {
	err := a.collection.FindOneAndUpdate(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"is_deleted": true}})
	return err.Err()
}
