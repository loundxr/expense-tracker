package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_postgres_pool "github.com/loundxr/expense-tracker/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) PatchUser(ctx context.Context, id int64, u domain.User) (domain.User, error) {
	const op = "users.repository.postgres"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	UPDATE expense_tracker.users
	SET
		email=$1,
		password_hash=$2,
		role=$3,
		version=version+1
	WHERE
		id=$4 AND version=$5
	RETURNING
		id, version, email, created_at, role`

	var um UserModel

	row := r.pool.QueryRow(ctx, query, u.Email, u.PasswordHash, u.Role, u.ID, u.Version)
	err := row.Scan(
		&um.ID,
		&um.Version,
		&um.Email,
		&um.CreatedAt,
		&um.Role,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNotFound) {
			return domain.User{}, fmt.Errorf("%s: version mismatch: %w", op, core_errors.ErrConflict)
		}
		if errors.Is(err, core_postgres_pool.ErrViolatesUniqueConstraint) {
			return domain.User{}, fmt.Errorf(
				"%s: scan from returned row: %w",
				op,
				core_errors.ErrAlreadyExists,
			)
		}
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	uDomain := userDomainFromModel(um)
	return uDomain, nil
}
