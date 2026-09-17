package repository

import (
	"context"
	"fmt"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

func (r *UsersAuthRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (domain.User, error) {
	const op = "repository.postgres.auth.GetUserByEmail"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	SELECT id, version, email, password_hash, created_at, role
	FROM expense_tracker.users
	WHERE email = $1`

	row := r.pool.QueryRow(ctx, query, email)

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
