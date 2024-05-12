package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/HeadGardener/coursework/internal/config"
	"github.com/HeadGardener/coursework/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	SecretKey       []byte
	AccessTokenTTL  time.Duration
	InitialLen      int
	RefreshTokenTTL time.Duration
}

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

func NewTokenManager(conf *config.TokensConfig) *TokenManager {
	return &TokenManager{
		SecretKey:       []byte(conf.SecretKey),
		AccessTokenTTL:  conf.AccessTokenTTL,
		InitialLen:      conf.InitialLen,
		RefreshTokenTTL: conf.RefreshTokenTTL,
	}
}

func (tm *TokenManager) GenerateAccessToken(userID string, role models.UserRole, age uint8) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		userID,
		role,
		age,
	})

	return token.SignedString(tm.SecretKey)
}

func (tm *TokenManager) ParseAccessToken(accessToken string) (UserAttributes, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return tm.SecretKey, nil
	})
	if err != nil {
		return UserAttributes{}, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return UserAttributes{}, errors.New("token claims are not of type *tokenClaims")
	}

	userAttributes := UserAttributes{
		ID:   claims.UserID,
		Role: claims.Role,
		Age:  claims.Age,
	}

	return userAttributes, nil
}

func (tm *TokenManager) ParseAccessTokenWithoutExpirationTime(accessToken string) (UserAttributes, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return tm.SecretKey, nil
	})
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return UserAttributes{}, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return UserAttributes{}, errors.New("token claims are not of type *tokenClaims")
	}

	userAttributes := UserAttributes{
		ID:   claims.UserID,
		Role: claims.Role,
		Age:  claims.Age,
	}

	return userAttributes, nil
}

func (tm *TokenManager) GenerateRefreshToken() (string, error) {
	b := make([]byte, tm.InitialLen)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	token := base64.StdEncoding.EncodeToString(b)

	return token, nil
}

func (tm *TokenManager) GetRefreshTokenTTL() time.Duration {
	return tm.RefreshTokenTTL
}
