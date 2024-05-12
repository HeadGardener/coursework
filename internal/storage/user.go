package storage

import (
	"context"

	"github.com/HeadGardener/coursework/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserStorage struct {
	db *sqlx.DB
}

func NewUserStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{db: db}
}

func (s *UserStorage) Create(ctx context.Context, user *models.User) (string, error) {
	if _, err := s.db.ExecContext(ctx, `insert into users (id, username, name, role, age, password_hash)
												values($1,$2,$3,$4,$5,$6)`,
		user.ID,
		user.Username,
		user.Name,
		user.Role,
		user.Age,
		user.PasswordHash); err != nil {
		return "", err
	}

	return user.ID, nil
}

func (s *UserStorage) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User

	if err := s.db.GetContext(ctx, &user, `select * from users where username=$1`, username); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserStorage) GetByID(ctx context.Context, userID string) (*models.User, error) {
	var user models.User

	if err := s.db.GetContext(ctx, &user, `select * from users where id=$1`, userID); err != nil {
		return nil, err
	}

	return &user, nil
}
