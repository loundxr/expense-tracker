package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

func (r *UsersRepository) GetUsers(ctx context.Context, limit, offset *int) ([]domain.User, error) {
	const op = "users.repository.postgres.GetUsers"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	SELECT id, version, email, password_hash, created_at, role
	FROM expense_tracker.users
	ORDER BY id ASC
	LIMIT $1
	OFFSET $2;`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: select users: %w", op, err)
	}
	defer rows.Close()

	var userModels []UserModel
	for rows.Next() {
		var um UserModel
		err := rows.Scan(
			&um.ID,
			&um.Version,
			&um.Email,
			&um.PasswordHash,
			&um.CreatedAt,
			&um.Role,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: scan row: %w", op, err)
		}
		userModels = append(userModels, um)
	}
	domains := userDomainsFromModels(userModels)
	return domains, nil
}
