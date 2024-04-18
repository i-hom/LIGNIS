package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Sale struct {
	AgentId    primitive.ObjectID `bson:"agent_id"`
	CustomerId primitive.ObjectID `bson:"customer_id"`
	SalesmanId primitive.ObjectID `bson:"salesman_id"`
	Cart       []ShortProduct     `bson:"cart"`
	Quantity   int32              `bson:"quantity"`
	TotalPrice float64            `bson:"total_price"`
}

type SaleWithID struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Sale `bson:",inline"`
}
