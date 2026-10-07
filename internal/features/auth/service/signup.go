package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
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
	var userDomain domain.User

	err = s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		var err error
		userDomain, err = s.usersRepository.CreateUser(txCtx, uninitUser)
		if err != nil {
			return err
		}

		defaultAccount := domain.NewUninitializedAccount("Personal", userDomain.ID)
		_, err = s.accountCreator.CreateAccount(txCtx, defaultAccount)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, core_errors.ErrAlreadyExists) {
			return domain.User{}, fmt.Errorf("%s: %w", op, err)
		}
		return domain.User{}, fmt.Errorf("failed to create user: %s: %w", op, err)
	}

	s.logger.Info(
		"user registered successfully with default 'Personal' account",
		slog.Int64("user_id", userDomain.ID),
		slog.String("email", userDomain.Email),
	)

	return userDomain, nil
}
