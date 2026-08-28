package auth_service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	"golang.org/x/crypto/bcrypt"
)

func (s *UsersAuthService) SignUp(ctx context.Context, email, password string) (domain.User, error) {
	const op = "users.service.auth.SignUp"

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(
			"failed to generate password hash",
			slog.String("error", err.Error()),
			slog.String("op", op),
		)
		return domain.User{}, fmt.Errorf("failed to generate password hash: %s: %w", op, err)
	}

	uninitUser := domain.NewUninitializedUser(email, string(passHash))

	userDomain, err := s.usersRepository.CreateUser(ctx, uninitUser)
	if err != nil {
		// TODO: check for ErrUserAlreadyExists
		s.logger.Error(
			"failed to create user",
			slog.String("error", err.Error()),
			slog.String("op", op),
		)
		return domain.User{}, fmt.Errorf("failed to create user: %s: %w", op, err)
	}

	// TODO: create account for the user by default

	return userDomain, nil
}
