package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/HeadGardener/coursework/internal/models"
	"github.com/redis/go-redis/v9"
)

type TokenStorage struct {
	rdb *redis.Client
}

func NewTokenStorage(rdb *redis.Client) *TokenStorage {
	return &TokenStorage{rdb: rdb}
}

func (s *TokenStorage) Add(ctx context.Context, session models.Session, ttl time.Duration) error {
	s.rdb.Del(ctx, session.UserID)

	b, err := models.MarshalSession(session)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	err = s.rdb.Set(ctx, session.UserID, b, ttl).Err()
	if err != nil {
		return fmt.Errorf("unable to store token: %w", err)
	}

	return nil
}

func (s *TokenStorage) Get(ctx context.Context, userID string) (models.Session, error) {
	b, err := s.rdb.Get(ctx, userID).Bytes()
	if err != nil {
		return models.Session{}, fmt.Errorf("user session doesn't exist: %w", err)
	}

	session, err := models.UnmarshalSession(b)
	if err != nil {
		return models.Session{}, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

func (s *TokenStorage) Delete(ctx context.Context, userID string) error {
	if err := s.rdb.Del(ctx, userID).Err(); err != nil {
		return err
	}

	return nil
}
