package auth

import (
	"errors"
	"time"

	"github.com/HeadGardener/coursework/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenTTL  = time.Hour
	secretKey = "qazwsxedcrfvtgbyhnujm"
)

type UserAttributes struct {
	ID   string          `json:"id"`
	Role models.UserRole `json:"user_role"`
	Age  uint8           `json:"age"`
}

type tokenClaims struct {
	jwt.RegisteredClaims
	UserID string          `json:"id"`
	Role   models.UserRole `json:"user_role"`
	Age    uint8           `json:"age"`
}

func GenerateToken(userID string, role models.UserRole, age uint8) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		userID,
		role,
		age,
	})

	return token.SignedString([]byte(secretKey))
}

func ParseToken(accessToken string) (UserAttributes, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return UserAttributes{}, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return UserAttributes{}, errors.New("token claims are not of type *tokenClaims")
	}

	return UserAttributes{
		ID:   claims.UserID,
		Role: claims.Role,
		Age:  claims.Age,
	}, nil
}
