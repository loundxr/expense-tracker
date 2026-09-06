package accounts_service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_ctx "github.com/loundxr/expense-tracker/internal/core/transport/http/context"
)

func (s *AccountsService) ShareAccess(ctx context.Context, id int, email string) error {
	const op = "accounts.service.ShareAccess"
	uid := core_ctx.GetUserID(ctx)
	role := core_ctx.GetUserRole(ctx)

	if err := s.accountsRepository.Exists(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if role != domain.RoleAdmin {
		isOwner, err := s.accountsRepository.IsOwner(ctx, uid, id)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if !isOwner {
			return fmt.Errorf("%s: only owner can share access to the account: %w", op, core_errors.ErrForbidden)
		}
	}

	targetUser, err := s.userProvider.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if targetUser.ID == uid {
		return fmt.Errorf("%s: you cannot grant access to yourself: %w", op, core_errors.ErrConflict)
	}

	if err := s.accountsRepository.GrantAccess(ctx, id, targetUser.ID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.logger.Info(
		"granted access",
		slog.Int("actor_id", uid),
		slog.Int("target_user_id", targetUser.ID),
		slog.Int("account_id", id),
	)

	return nil
}
