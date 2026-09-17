package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_postgres_pool "github.com/loundxr/expense-tracker/internal/core/repository/postgres/pool"
)

func (r *UsersAuthRepository) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	const op = "repository.postgres.Auth.CreateUser"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	INSERT INTO expense_tracker.users (email, password_hash, created_at)
	VALUES($1, $2, $3)
	RETURNING id, version, email, password_hash, created_at, role;`

	row := r.pool.QueryRow(ctx, query, user.Email, user.PasswordHash, user.CreatedAt)

	var um UserModel
	err := row.Scan(
		&um.ID,
		&um.Version,
		&um.Email,
		&um.PasswordHash,
		&um.CreatedAt,
		&um.Role,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesUniqueConstraint) {
			return domain.User{}, fmt.Errorf(
				"%s: scan from returned row: %w",
				op,
				core_errors.ErrAlreadyExists,
			)
		}
		return domain.User{}, fmt.Errorf("%s: scan from returned row: %w", op, err)
	}

	return domain.NewUser(
		um.ID,
		um.Version,
		um.Email,
		um.PasswordHash,
		um.CreatedAt,
		um.Role,
	), nil
}
