package message

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog"
)

type Message struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Content   string    `db:"content"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	ExpiresAt time.Time `db:"expires_at"`
}

type Repository interface {
	CreateMessage(ctx context.Context, msg *Message) error
	GetMessages(ctx context.Context, limit int) ([]Message, error)
	DeleteExpiredMessages(ctx context.Context) error
}

type Service struct {
	repo         Repository
	log          zerolog.Logger
	messageTTL   time.Duration
	cleanupTimer *time.Ticker
}

func NewService(repo Repository, log zerolog.Logger, messageTTL time.Duration) *Service {
	return &Service{
		repo:       repo,
		log:        log,
		messageTTL: messageTTL,
	}
}

func (s *Service) StartCleanupRoutine(ctx context.Context, interval time.Duration) {
	s.cleanupTimer = time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-s.cleanupTimer.C:
				if err := s.repo.DeleteExpiredMessages(ctx); err != nil {
					s.log.Error().Err(err).Msg("Failed to delete expired messages")
				}
			case <-ctx.Done():
				s.cleanupTimer.Stop()
				return
			}
		}
	}()
}

func (s *Service) PostMessage(ctx context.Context, userID int, content string) (*Message, error) {
	if content == "" {
		return nil, errors.New("message content cannot be empty")
	}

	msg := &Message{
		UserID:    userID,
		Content:   content,
		Status:    "active",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(s.messageTTL),
	}

	if err := s.repo.CreateMessage(ctx, msg); err != nil {
		s.log.Error().Err(err).Msg("Failed to create message")
		return nil, errors.New("failed to create message")
	}

	return msg, nil
}

func (s *Service) GetRecentMessages(ctx context.Context, limit int) ([]Message, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.GetMessages(ctx, limit)
}