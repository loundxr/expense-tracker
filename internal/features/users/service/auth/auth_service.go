package users_service_auth

import (
	"context"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_logger "github.com/loundxr/expense-tracker/internal/core/logger"
)

type UsersAuthService struct {
	usersRepository UsersAuthRepository
	logger          *core_logger.Logger
}

type UsersAuthRepository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
}

func NewUsersAuthService(ur UsersAuthRepository, l *core_logger.Logger) *UsersAuthService {
	return &UsersAuthService{
		usersRepository: ur,
		logger:          l,
	}
}
