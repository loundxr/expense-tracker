package users_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_postgres_pool "github.com/loundxr/expense-tracker/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) GetUserByID(ctx context.Context, id int) (domain.User, error) {
	const op = "users.repository.postgres.GetUser"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	SELECT id, version, email, password_hash, created_at, role
	FROM expense_tracker.users
	WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)

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
		if errors.Is(err, core_postgres_pool.ErrNotFound) {
			return domain.User{}, fmt.Errorf("%w", core_errors.ErrNotFound)
		}
		return domain.User{}, fmt.Errorf("%s: scan row: %w", op, err)
	}
	return userDomainFromModel(um), nil
}

func (r *UsersRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	const op = "users.repository.postgres.GetUserByEmail"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	SELECT id, version, email, created_at, role
	FROM expense_tracker.users
	WHERE email = $1`

	row := r.pool.QueryRow(ctx, query, email)

	var um UserModel
	err := row.Scan(
		&um.ID,
		&um.Version,
		&um.Email,
		&um.CreatedAt,
		&um.Role,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNotFound) {
			return domain.User{}, fmt.Errorf("%s: %w", op, core_errors.ErrNotFound)
		}
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	userDomain := userDomainFromModel(um)
	return userDomain, nil
}
