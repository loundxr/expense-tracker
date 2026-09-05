package accounts_service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_ctx "github.com/loundxr/expense-tracker/internal/core/transport/http/context"
)

func (s *AccountsService) GetAccounts(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.Account, error) {
	const op = "accounts.service.GetAccounts"

	uid := core_ctx.GetUserID(ctx)
	role := core_ctx.GetUserRole(ctx)

	var (
		accounts []domain.Account
		err      error
	)

	if role == domain.RoleAdmin {
		s.logger.Info(
			"admin fetching all system accounts",
			slog.Int("id", uid),
		)
		accounts, err = s.accountsRepository.GetAllAccounts(ctx, limit, offset)
	} else {
		s.logger.Info(
			"user fetching available accounts",
			slog.Int("id", uid),
		)
		accounts, err = s.accountsRepository.GetAccountByUserID(ctx, uid, limit, offset)
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return accounts, nil
}
