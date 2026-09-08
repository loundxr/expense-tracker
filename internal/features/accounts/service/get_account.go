package accounts_service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_ctx "github.com/loundxr/expense-tracker/internal/core/transport/http/context"
)

func (s *AccountsService) GetAccount(ctx context.Context, accountID int64) (domain.Account, error) {
	const op = "accounts.service.GetAccount"

	uid := core_ctx.GetUserID(ctx)
	role := core_ctx.GetUserRole(ctx)

	account, err := s.accountsRepository.GetAccountByID(ctx, accountID)
	if err != nil {
		return domain.Account{}, fmt.Errorf("%s: %w", op, err)
	}

	if role != domain.RoleAdmin {
		hasAccess, err := s.accountsRepository.HasAccess(ctx, uid, accountID)
		if err != nil {
			return domain.Account{}, fmt.Errorf("%s: %w", op, err)
		}

		if !hasAccess {
			return domain.Account{}, fmt.Errorf("%s: %w", op, core_errors.ErrForbidden)
		}
	}

	s.logger.Info(
		"account info fetched",
		slog.Int64("actor_id", uid),
		slog.Int64("account_id", accountID),
	)
	return account, nil
}
