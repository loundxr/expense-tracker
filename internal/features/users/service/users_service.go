package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

type Cache interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
}

type UsersService struct {
	usersRepository UsersRepository
	revokeStore     Cache
	logger          *slog.Logger
}

type UsersRepository interface {
	GetUsers(ctx context.Context, limit *int, offset *int) ([]domain.User, error)
	GetUserByID(ctx context.Context, id int64) (domain.User, error)
	DeleteUser(ctx context.Context, id int64) error
	PatchUser(ctx context.Context, id int64, user domain.User) (domain.User, error)
}

func NewUsersService(ur UsersRepository, rs Cache, l *slog.Logger) *UsersService {
	return &UsersService{
		usersRepository: ur,
		revokeStore:     rs,
		logger:          l,
	}
}
