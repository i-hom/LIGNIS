package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Defect struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Defects   []DefectProduct    `bson:"defects"`
	CreatedBy primitive.ObjectID `bson:"created_by"`
}

type DefectProduct struct {
	ProductID primitive.ObjectID `bson:"product_id"`
	Quantity  uint32             `bson:"quantity"`
	Remark    string             `bson:"remark"`
}
