package repository

import (
	"context"

	"lignis/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DefectRepo struct {
	collection *mongo.Collection
}

func NewDefectRepo(collection *mongo.Collection) *DefectRepo {
	return &DefectRepo{
		collection: collection,
	}
}

func (r DefectRepo) Create(defect *model.Defect) (primitive.ObjectID, error) {

	res, err := r.collection.InsertOne(context.TODO(), defect)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (r DefectRepo) GetByPatter(pattern string, limit, page int64) ([]model.Defect, int64, error) {

	var defects []model.Defect
	filter := bson.M{
		"defects.remark": bson.M{"$regex": pattern, "$options": "i"},
	}
	cursor, err := r.collection.Find(
		context.TODO(),
		filter,
		options.Find().SetSort(bson.D{{Key: "_id", Value: 1}}).SetLimit(limit).SetSkip((page-1)*limit))

	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.TODO())

	count, err := r.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, 0, err
	}
	for cursor.Next(context.TODO()) {
		var defect model.Defect
		err = cursor.Decode(&defect)
		if err != nil {
			return nil, 0, err
		}
		defects = append(defects, defect)
	}
	return defects, count, nil
}

func (r DefectRepo) GetByID(id primitive.ObjectID) (*model.Defect, error) {
	var defect model.Defect
	err := r.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&defect)
	if err != nil {
		return nil, err
	}
	return &defect, nil
}

func (r DefectRepo) Delete(id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}
