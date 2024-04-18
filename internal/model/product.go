package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	Name      string  `bson:"name"`
	Code      string  `bson:"code"`
	Quantity  uint32  `bson:"quantity"`
	CostPrice float64 `bson:"price"`
	SellPrice float64 `bson:"sell_price"`
}

type ProductWithID struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	Product `bson:",inline"`
}

type ShortProduct struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Quantity uint32             `bson:"quantity"`
	Price    float64            `bson:"price"`
}
