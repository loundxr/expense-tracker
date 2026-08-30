package auth_service

import (
	"context"
	"log/slog"

	core_jwt "github.com/loundxr/expense-tracker/internal/core/auth/jwt"
	"github.com/loundxr/expense-tracker/internal/core/domain"
)

type UsersAuthService struct {
	usersRepository UsersAuthRepository
	logger          *slog.Logger
	cfg             core_jwt.JWTConfig
}

type UsersAuthRepository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
}

func NewUsersAuthService(ur UsersAuthRepository, l *slog.Logger, jwtCfg core_jwt.JWTConfig) *UsersAuthService {
	return &UsersAuthService{
		usersRepository: ur,
		logger:          l,
		cfg:             jwtCfg,
	}
}
