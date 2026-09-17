package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_ctx "github.com/loundxr/expense-tracker/internal/core/transport/http/context"
)

func (s *AccountsService) DeleteAccount(ctx context.Context, accountID int64) error {
	const op = "accounts.service.DeleteAccount"
	uid := core_ctx.GetUserID(ctx)
	role := core_ctx.GetUserRole(ctx)

	if err := s.accountsRepository.Exists(ctx, accountID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if role != domain.RoleAdmin {
		isOwner, err := s.accountsRepository.IsOwner(ctx, uid, accountID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if !isOwner {
			return fmt.Errorf("%s: %w", op, core_errors.ErrForbidden)
		}
	}

	err := s.accountsRepository.DeleteAccount(ctx, accountID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.logger.Info(
		"account deleted",
		slog.Int64("actor_id", uid),
		slog.Int64("account_id", accountID),
	)

	return nil
}
