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
	GetAccountByUserID(ctx context.Context, uid int, limit *int, offset *int) ([]domain.Account, error)
}

func NewAccountsService(ar AccountsRepository, l *slog.Logger) *AccountsService {
	return &AccountsService{
		accountsRepository: ar,
		logger:             l,
	}
}
