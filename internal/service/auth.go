package service

import (
	"context"
	"errors"

	"github.com/HeadGardener/coursework/internal/lib/auth"

	"github.com/HeadGardener/coursework/internal/lib/hash"
	"github.com/HeadGardener/coursework/internal/models"
	"github.com/google/uuid"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
)

type TokenStorage interface {
	Add(ctx context.Context, userID, token string) error
	Check(ctx context.Context, userID, token string) error
	Delete(ctx context.Context, userID string) error
}

type UserStorage interface {
	Create(ctx context.Context, user *models.User) (string, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
}

type AuthService struct {
	tokenStorage TokenStorage
	userStorage  UserStorage
}

func NewAuthService(tokenStorage TokenStorage, userStorage UserStorage) *AuthService {
	return &AuthService{
		tokenStorage: tokenStorage,
		userStorage:  userStorage,
	}
}

func (s *AuthService) SignUp(ctx context.Context, username, name string, age int, password string) (string, error) {
	user := &models.User{
		ID:           uuid.NewString(),
		Username:     username,
		Name:         name,
		Role:         models.RoleUser,
		Age:          uint8(age),
		PasswordHash: hash.GetPasswordHash(password),
	}

	return s.userStorage.Create(ctx, user)
}

func (s *AuthService) SignIn(ctx context.Context, username, password string) (string, error) {
	user, err := s.userStorage.GetByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	if !hash.CheckPassword([]byte(user.PasswordHash), password) {
		return "", ErrInvalidPassword
	}

	token, err := auth.GenerateToken(user.ID, user.Role, user.Age)
	if err != nil {
		return "", err
	}

	if err = s.tokenStorage.Add(ctx, user.ID, token); err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Check(ctx context.Context, userID, token string) error {
	return s.tokenStorage.Check(ctx, userID, token)
}

func (s *AuthService) LogOut(ctx context.Context, userID string) error {
	return s.tokenStorage.Delete(ctx, userID)
}
