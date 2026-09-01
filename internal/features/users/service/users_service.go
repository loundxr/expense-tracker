package users_service

import (
	"context"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

type UsersService struct {
	usersRepository UsersRepository
	logger          *slog.Logger
}

type UsersRepository interface {
	GetUsers(ctx context.Context, limit *int, offset *int) ([]domain.User, error)
	GetUserByID(ctx context.Context, id int) (domain.User, error)
	DeleteUser(ctx context.Context, id int) error
	PatchUser(ctx context.Context, id int, user domain.User) (domain.User, error)
}

func NewUsersService(ur UsersRepository, l *slog.Logger) *UsersService {
	return &UsersService{
		usersRepository: ur,
		logger:          l,
	}
}
