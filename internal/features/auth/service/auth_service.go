package auth_service

import (
	"context"
	"log/slog"

	core_jwt "github.com/loundxr/expense-tracker/internal/core/auth/jwt"
	"github.com/loundxr/expense-tracker/internal/core/domain"
)

type AccountCreator interface {
	CreateAccount(ctx context.Context, account domain.Account) (domain.Account, error)
}

type UsersAuthService struct {
	usersRepository UsersAuthRepository
	accountCreator  AccountCreator
	logger          *slog.Logger
	cfg             core_jwt.JWTConfig
}

type UsersAuthRepository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
}

func NewUsersAuthService(
	ur UsersAuthRepository,
	ac AccountCreator,
	l *slog.Logger,
	jwtCfg core_jwt.JWTConfig,
) *UsersAuthService {
	return &UsersAuthService{
		usersRepository: ur,
		accountCreator:  ac,
		logger:          l,
		cfg:             jwtCfg,
	}
}
