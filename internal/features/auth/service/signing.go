package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	core_jwt "github.com/loundxr/expense-tracker/internal/core/auth/jwt"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	"golang.org/x/crypto/bcrypt"
)

func (s *UsersAuthService) SignIn(ctx context.Context, email, password string) (string, error) {
	const op = "users.service.auth.SignIn"

	user, err := s.usersRepository.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return "", fmt.Errorf("%s: %w", op, core_errors.ErrUnauthorized)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, core_errors.ErrUnauthorized)
	}

	token, err := core_jwt.NewToken(user, s.cfg.Secret, s.cfg.TTL)
	if err != nil {
		s.logger.Error("failed to generate token", slog.String("op", op), slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	s.logger.Info("user logged in", slog.Int64("user_id", user.ID))

	return token, nil
}
