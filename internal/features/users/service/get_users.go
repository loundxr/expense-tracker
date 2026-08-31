package users_service

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

	if limit != nil && *limit <= 0 {
		return nil, fmt.Errorf(
			"limit must be positive: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if limit != nil && *limit > 100 {
		*limit = 100
	}

	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf(
			"offset must be non-negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	users, err := s.usersRepository.GetUsers(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	s.logger.Debug("admin accessed the users list", slog.Int("id", core_ctx.GetUserID(ctx)))
	return users, nil
}
