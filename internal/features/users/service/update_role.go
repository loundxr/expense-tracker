package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_ctx "github.com/loundxr/expense-tracker/internal/core/transport/http/context"
)

func (s *UsersService) UpdateUserRole(ctx context.Context, id int64, role string) (domain.User, error) {
	const op = "users.service.UpdateUserRole"
	currID := core_ctx.GetUserID(ctx)
	currRole := core_ctx.GetUserRole(ctx)

	if currRole != domain.RoleAdmin {
		return domain.User{}, fmt.Errorf("%s: %w", op, core_errors.ErrForbidden)
	}

	user, err := s.usersRepository.GetUserByID(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	prevRole := user.Role
	if role == prevRole {
		return domain.User{}, fmt.Errorf("%s: role matches the current user's role: %w", op, core_errors.ErrInvalidArgument)
	}
	user.Role = role

	updatedUser, err := s.usersRepository.PatchUser(ctx, id, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	revokeKey := fmt.Sprintf("auth:revoked_at:%d", id)
	if err := s.revokeStore.Set(ctx, revokeKey, time.Now().Unix(), 24*time.Hour); err != nil {
		s.logger.Error(
			"failed to set token revocation timestamp in cache",
			slog.Int64("target_id", id),
			slog.String("error", err.Error()),
		)
	}

	s.logger.Warn(
		"user role updated",
		slog.Int64("actor_id", currID),
		slog.Int64("target_id", id),
		slog.String("previous_role", prevRole),
		slog.String("new_role", role),
	)

	return updatedUser, nil
}
