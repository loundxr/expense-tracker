package users_service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_ctx "github.com/loundxr/expense-tracker/internal/core/transport/http/context"
)

func (s *UsersService) DeleteUser(ctx context.Context, id int) error {
	const op = "users.service.DeleteUser"

	currUID := core_ctx.GetUserID(ctx)
	currRole := core_ctx.GetUserRole(ctx)

	if currRole != domain.RoleAdmin && currUID != id {
		return fmt.Errorf("%s: %w", op, core_errors.ErrForbidden)
	}

	if err := s.usersRepository.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.logger.Info(
		"user profile deleted",
		slog.Int("actor_id", currUID),
		slog.String("actor_role", currRole),
		slog.Int("target_id", id),
	)

	return nil
}
