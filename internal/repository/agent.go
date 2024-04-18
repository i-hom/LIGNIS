package repository

import (
	"context"
	"lignis/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AgentRepo struct {
	collection *mongo.Collection
}

func NewAgentRepo(collection *mongo.Collection) *AgentRepo {
	return &AgentRepo{
		collection: collection,
	}
}

func (a *AgentRepo) Create(agent *model.Agent) (primitive.ObjectID, error) {
	res, err := a.collection.InsertOne(context.TODO(), agent)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}
