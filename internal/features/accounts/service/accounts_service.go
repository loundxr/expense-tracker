package accounts_service

import (
	"context"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

type AccountsService struct {
	accountsRepository AccountsRepository
	logger             *slog.Logger
}

type AccountsRepository interface {
	CreateAccount(ctx context.Context, account domain.Account) (domain.Account, error)
	GetAllAccounts(ctx context.Context, limit *int, offset *int) ([]domain.Account, error)
	GetAccountsByUserID(ctx context.Context, uid int, limit *int, offset *int) ([]domain.Account, error)
	GetAccountByID(ctx context.Context, id int) (domain.Account, error)
	DeleteAccount(ctx context.Context, accountID int) error
	HasAccess(ctx context.Context, uid int, accountID int) (bool, error)
	IsOwner(ctx context.Context, uid int, accountID int) (bool, error)
	Exists(ctx context.Context, accountID int) error
}

func NewAccountsService(ar AccountsRepository, l *slog.Logger) *AccountsService {
	return &AccountsService{
		accountsRepository: ar,
		logger:             l,
	}
}
