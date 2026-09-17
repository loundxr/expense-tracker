package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

func (s *AccountsService) CreateAccount(ctx context.Context, account domain.Account) (domain.Account, error) {
	const op = "accounts.service.CreateAccount"

	acc, err := s.accountsRepository.CreateAccount(ctx, account)
	if err != nil {
		return domain.Account{}, fmt.Errorf("%s: %w", op, err)
	}

	s.logger.Info(
		"account created",
		slog.Int64("account_id", acc.ID),
		slog.Int64("owner_id", acc.UserID),
	)
	return acc, nil
}
