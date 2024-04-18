package model

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Claims struct {
	UserID primitive.ObjectID `json:"user_id"`
	Role   string             `json:"role"`
	jwt.StandardClaims
}
