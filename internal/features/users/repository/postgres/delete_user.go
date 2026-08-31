package users_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
)

func (r *UsersRepository) DeleteUser(ctx context.Context, id int) error {
	const op = "users.repository.postgres.DeleteUser"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	DELETE FROM expense_tracker.users
	WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user with id=%d: %w", id, core_errors.ErrNotFound)
	}
	return nil
}
