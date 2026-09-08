package users_service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_ctx "github.com/loundxr/expense-tracker/internal/core/transport/http/context"
	"golang.org/x/crypto/bcrypt"
)

func (s *UsersService) PatchUser(ctx context.Context, id int64, patch domain.UserPatch) (domain.User, error) {
	const op = "users.service.PatchUser"
	currID := core_ctx.GetUserID(ctx)
	currRole := core_ctx.GetUserRole(ctx)

	if currRole != domain.RoleAdmin && currID != id {
		return domain.User{}, fmt.Errorf("%s: %w", op, core_errors.ErrForbidden)
	}

	user, err := s.usersRepository.GetUserByID(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: failed to get user: %w", op, err)
	}

	var patchedFields []string

	if patch.Password.Set {
		hash, err := bcrypt.GenerateFromPassword([]byte(*patch.Password.Val), bcrypt.DefaultCost)
		if err != nil {
			return domain.User{}, fmt.Errorf("%s: hash password: %w", op, err)
		}
		user.PasswordHash = string(hash)
		patchedFields = append(patchedFields, "password")
	}
	if patch.Email.Set {
		patchedFields = append(patchedFields, "email")
	}

	if err := user.ApplyPatch(patch); err != nil {
		return domain.User{}, fmt.Errorf("%s: apply patch: %w", op, err)
	}

	patchedUser, err := s.usersRepository.PatchUser(ctx, id, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	s.logger.Info(
		"user profile patched",
		slog.Int64("actor_id", currID),
		slog.Int64("target_id", id),
		slog.Any("fields", patchedFields),
	)

	return patchedUser, nil
}
