package service

import (
	"context"
	"errors"
	"time"

	"github.com/HeadGardener/coursework/internal/lib/auth"

	"github.com/HeadGardener/coursework/internal/lib/hash"
	"github.com/HeadGardener/coursework/internal/models"
	"github.com/google/uuid"
)

var (
	ErrInvalidPassword     = errors.New("invalid password")
	ErrNotSameRefreshToken = errors.New("invalid refresh token")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
	ErrInvalidSession      = errors.New("tokens are not connected: invalid access token session")
)

type TokenManager interface {
	GenerateAccessToken(userID string, role models.UserRole, age uint8) (string, error)
	ParseAccessToken(accessToken string) (auth.UserAttributes, error)
	ParseAccessTokenWithoutExpirationTime(accessToken string) (auth.UserAttributes, error)
	GenerateRefreshToken() (string, error)
	GetRefreshTokenTTL() time.Duration
}

type SessionStorage interface {
	Add(ctx context.Context, session models.Session, ttl time.Duration) error
	Get(ctx context.Context, userID string) (models.Session, error)
	Delete(ctx context.Context, userID string) error
}

type UserStorage interface {
	Create(ctx context.Context, user *models.User) (string, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByID(ctx context.Context, userID string) (*models.User, error)
}

type AuthService struct {
	tokenManager   TokenManager
	sessionStorage SessionStorage
	userStorage    UserStorage
}

func NewAuthService(tokenManager TokenManager, tokenStorage SessionStorage, userStorage UserStorage) *AuthService {
	return &AuthService{
		tokenManager:   tokenManager,
		sessionStorage: tokenStorage,
		userStorage:    userStorage,
	}
}

func (s *AuthService) SignUp(ctx context.Context, username, name string, age int, password string) (string, error) {
	user := &models.User{
		ID:           uuid.NewString(),
		Username:     username,
		Name:         name,
		Role:         models.RoleUser,
		Age:          uint8(age),
		PasswordHash: hash.GetStringHash(password),
	}

	return s.userStorage.Create(ctx, user)
}

func (s *AuthService) SignIn(ctx context.Context, username, password string) (models.Tokens, error) {
	user, err := s.userStorage.GetByUsername(ctx, username)
	if err != nil {
		return models.Tokens{}, err
	}

	if !hash.CompareHashAndString([]byte(user.PasswordHash), password) {
		return models.Tokens{}, ErrInvalidPassword
	}

	tokens, err := s.createSession(ctx, user.ID, user.Role, user.Age)
	if err != nil {
		return models.Tokens{}, err
	}

	return tokens, nil
}

func (s *AuthService) ParseAccessToken(token string) (auth.UserAttributes, error) {
	return s.tokenManager.ParseAccessToken(token)
}

func (s *AuthService) Refresh(ctx context.Context, accessToken, refreshToken string) (models.Tokens, error) {
	userAttr, err := s.tokenManager.ParseAccessTokenWithoutExpirationTime(accessToken)
	if err != nil {
		return models.Tokens{}, err
	}

	session, err := s.sessionStorage.Get(ctx, userAttr.ID)
	if err != nil {
		return models.Tokens{}, err
	}

	if !hash.CompareHashAndString([]byte(session.RefreshToken), refreshToken) {
		return models.Tokens{}, ErrNotSameRefreshToken
	}

	if session.ExpiresAt.Before(time.Now()) {
		return models.Tokens{}, ErrRefreshTokenExpired
	}

	user, err := s.userStorage.GetByID(ctx, session.UserID)
	if err != nil {
		return models.Tokens{}, err
	}

	tokens, err := s.createSession(ctx, user.ID, user.Role, user.Age)
	if err != nil {
		return models.Tokens{}, err
	}

	return tokens, nil
}

func (s *AuthService) LogOut(ctx context.Context, userID string) error {
	return s.sessionStorage.Delete(ctx, userID)
}

func (s *AuthService) createSession(ctx context.Context, userID string, role models.UserRole, age uint8) (models.Tokens, error) {
	var (
		tokens models.Tokens
		err    error
	)

	sessionID := uuid.NewString()
	tokens.AccessToken, err = s.tokenManager.GenerateAccessToken(userID, role, age)
	if err != nil {
		return models.Tokens{}, err
	}

	tokens.RefreshToken, err = s.tokenManager.GenerateRefreshToken()
	if err != nil {
		return models.Tokens{}, err
	}

	session := models.Session{
		ID:           sessionID,
		UserID:       userID,
		RefreshToken: hash.GetStringHash(tokens.RefreshToken),
		ExpiresAt:    time.Now().Add(s.tokenManager.GetRefreshTokenTTL()),
	}

	if err = s.sessionStorage.Add(ctx, session, s.tokenManager.GetRefreshTokenTTL()); err != nil {
		return models.Tokens{}, err
	}

	return tokens, nil
}
