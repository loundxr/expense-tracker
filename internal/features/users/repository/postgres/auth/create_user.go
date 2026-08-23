package users_repository_postgres_auth

import (
	"context"
	"fmt"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	users_postgres_repository "github.com/loundxr/expense-tracker/internal/features/users/repository/postgres"
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
	RETURNING id, version, email, password_hash, created_at;`

	row := r.pool.QueryRow(ctx, query, user.Email, user.PasswordHash, user.CreatedAt)

	var um users_postgres_repository.UserModel
	err := row.Scan(
		&um.ID,
		&um.Version,
		&um.Email,
		&um.PasswordHash,
		&um.CreatedAt,
	)

	if err != nil {
		return domain.User{}, fmt.Errorf("scan from returned row: %w", err)
	}

	userDomain := domain.NewUser(
		um.ID,
		um.Version,
		um.Email,
		um.PasswordHash,
		um.CreatedAt,
	)
	return userDomain, nil
}
