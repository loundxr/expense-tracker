package users_service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_ctx "github.com/loundxr/expense-tracker/internal/core/transport/http/context"
)

// TODO: implement redis 'iat' state for dynamically checking time user was updated
// to prevent non-admins from manipulating data
func (s *UsersService) UpdateUserRole(ctx context.Context, id int, role string) (domain.User, error) {
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
		return domain.User{}, fmt.Errorf("%s: %w", err)
	}

	s.logger.Warn(
		"user role updated",
		slog.Int("actor_id", currID),
		slog.Int("target_id", id),
		slog.String("previous_role", prevRole),
		slog.String("new_role", role),
	)

	return updatedUser, nil
}
