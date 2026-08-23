package users_service_auth

import (
	"context"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

type UsersAuthService struct {
	usersRepository UsersAuthRepository
	logger          *slog.Logger
}

type UsersAuthRepository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
}

func NewUsersAuthService(ur UsersAuthRepository, l *slog.Logger) *UsersAuthService {
	return &UsersAuthService{
		usersRepository: ur,
		logger:          l,
	}
}
