package accounts_service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_ctx "github.com/loundxr/expense-tracker/internal/core/transport/http/context"
)

func (s *AccountsService) RevokeAccess(ctx context.Context, accountID int, email string) error {
	const op = "accounts.service.RevokeAccess"
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
			return fmt.Errorf("%s: you can revoke access only from your account: %w", op, core_errors.ErrForbidden)
		}
	}

	targetUser, err := s.userProvider.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if targetUser.ID == uid {
		return fmt.Errorf("%s: you cannot revoke access from yourself: %w", op, core_errors.ErrConflict)
	}

	if err := s.accountsRepository.RevokeAccess(ctx, accountID, targetUser.ID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.logger.Info(
		"revoked access",
		slog.Int("actor_id", uid),
		slog.Int("target_id", targetUser.ID),
		slog.Int("account_id", accountID),
	)
	return nil
}
