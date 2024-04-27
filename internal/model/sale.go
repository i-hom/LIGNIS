package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Sale struct {
	AgentId      primitive.ObjectID `bson:"agent_id"`
	CustomerId   primitive.ObjectID `bson:"customer_id"`
	SalesmanId   primitive.ObjectID `bson:"salesman_id"`
	Cart         []ShortProduct     `bson:"cart"`
	TotalUZS     float64            `bson:"total_uzs"`
	TotalUSD     float64            `bson:"total_usd"`
	CurrencyCode string             `bson:"currency_code"`
	Deleted_At   primitive.DateTime `bson:"deleted_at,omitempty"`
}

type SaleWithID struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Sale `bson:",inline"`
}
