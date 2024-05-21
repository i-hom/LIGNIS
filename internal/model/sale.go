package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Sale struct {
	AgentId      primitive.ObjectID `bson:"agent_id,omitempty"`
	CustomerId   primitive.ObjectID `bson:"customer_id,omitempty"`
	SalesmanId   primitive.ObjectID `bson:"salesman_id"`
	Cart         []ShortProduct     `bson:"cart"`
	TotalUZS     int64              `bson:"total_uzs"`
	TotalUSD     float64            `bson:"total_usd"`
	CurrencyCode string             `bson:"currency_code"`
	Is_Deleted   bool               `bson:"deleted_at,omitempty"`
	Deleted_By   primitive.ObjectID `bson:"deleted_by,omitempty"`
}

type SaleWithID struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Sale `bson:",inline"`
	Date time.Time
}
