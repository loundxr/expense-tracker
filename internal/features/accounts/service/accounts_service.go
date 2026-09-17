package service

import (
	"context"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

type AccountsService struct {
	accountsRepository AccountsRepository
	userProvider       UserProvider
	logger             *slog.Logger
}

type AccountsRepository interface {
	CreateAccount(ctx context.Context, account domain.Account) (domain.Account, error)
	GetAllAccounts(ctx context.Context, limit *int, offset *int) ([]domain.Account, error)
	GetAccountsByUserID(ctx context.Context, uid int64, limit *int, offset *int) ([]domain.Account, error)
	GetAccountByID(ctx context.Context, id int64) (domain.Account, error)
	DeleteAccount(ctx context.Context, accountID int64) error
	PatchAccount(ctx context.Context, accountID int64, toPatch domain.Account) (domain.Account, error)
	GrantAccess(ctx context.Context, accountID int64, targetID int64) error
	RevokeAccess(ctx context.Context, accountID int64, targetID int64) error
	HasAccess(ctx context.Context, uid int64, accountID int64) (bool, error)
	IsOwner(ctx context.Context, uid int64, accountID int64) (bool, error)
	Exists(ctx context.Context, accountID int64) error
}

type UserProvider interface {
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
}

func NewAccountsService(ar AccountsRepository, up UserProvider, l *slog.Logger) *AccountsService {
	return &AccountsService{
		accountsRepository: ar,
		userProvider:       up,
		logger:             l,
	}
}
