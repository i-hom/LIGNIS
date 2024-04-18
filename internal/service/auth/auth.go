package auth

import (
	"context"
	"lignis/internal/model"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Auth struct {
	signKey string
	ttl     int
}

func NewAuth(config *Config) *Auth {
	return &Auth{
		signKey: config.SignKey,
		ttl:     config.TTL,
	}
}

func (a Auth) GenerateToken(user *model.UserWithID) (string, error) {
	expirationTime := time.Now().Add(time.Duration(a.ttl))
	claims := &model.Claims{
		UserID: user.ID,
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.signKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a Auth) ValidateAndParseToken(ctx context.Context, tokenString string) (*model.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.signKey), nil
	})
	if err != nil {
		return &model.Claims{}, err
	}
	claims, ok := token.Claims.(*model.Claims)
	if !ok || !token.Valid {
		return &model.Claims{}, err
	}

	return claims, nil
}
