package users_service

import (
	"context"
	"fmt"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_ctx "github.com/loundxr/expense-tracker/internal/core/transport/http/context"
)

func (s *UsersService) GetUser(ctx context.Context, id int) (domain.User, error) {
	const op = "users.service.GetUser"

	currUID := core_ctx.GetUserID(ctx)
	currRole := core_ctx.GetUserRole(ctx)

	if currRole != domain.RoleAdmin && currUID != id {
		return domain.User{}, fmt.Errorf("%s: %w", op, core_errors.ErrForbidden)
	}

	user, err := s.usersRepository.GetUserByID(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
