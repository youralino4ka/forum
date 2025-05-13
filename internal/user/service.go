package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/rs/zerolog"
)

type User struct {
	ID           int       `db:"id"`
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Role         string    `db:"role"`
	Status       string    `db:"status"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id int) error
}

type Service struct {
	repo Repository
	log  zerolog.Logger
}

func NewService(repo Repository, log zerolog.Logger) *Service {
	return &Service{repo: repo, log: log}
}

func (s *Service) Register(ctx context.Context, username, email, password string) (*User, error) {
	// Проверка существования пользователя
	_, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error().Err(err).Msg("Failed to check user existence")
		return nil, errors.New("failed to check user existence")
	} else if err == nil {
		return nil, errors.New("username already exists")
	}

	// Хеширование пароля (в реальном проекте используйте bcrypt)
	passwordHash := "hashed_" + password // Замените на реальное хеширование

	user := &User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         "user",
		Status:       "active",
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		s.log.Error().Err(err).Msg("Failed to create user")
		return nil, errors.New("failed to create user")
	}

	return user, nil
}

// Другие методы сервиса...