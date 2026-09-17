package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_ctx "github.com/loundxr/expense-tracker/internal/core/transport/http/context"
)

func (s *UsersService) GetUsers(ctx context.Context, limit *int, offset *int) ([]domain.User, error) {
	const op = "users.service.GetUsers"

	if core_ctx.GetUserRole(ctx) != domain.RoleAdmin {
		return nil, fmt.Errorf("%s: %w", op, core_errors.ErrForbidden)
	}

	users, err := s.usersRepository.GetUsers(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	s.logger.Info(
		"users list fetched by admin",
		slog.Int64("actor_id", core_ctx.GetUserID(ctx)),
	)
	return users, nil
}
