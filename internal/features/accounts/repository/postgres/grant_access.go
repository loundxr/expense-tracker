package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_postgres_pool "github.com/loundxr/expense-tracker/internal/core/repository/postgres/pool"
)

func (r *AccountsRepository) GrantAccess(ctx context.Context, accountID int64, targetID int64) error {
	const op = "accounts.repository.postgres.GrantAccess"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	INSERT INTO expense_tracker.account_users (account_id, user_id, created_at)
	VALUES($1, $2, $3);`

	_, err := r.pool.Exec(ctx, query, accountID, targetID, time.Now())
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesUniqueConstraint) {
			return fmt.Errorf("%s: user already has access: %w", op, core_errors.ErrAlreadyExists)
		}
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return fmt.Errorf("%s: %w", op, core_errors.ErrNotFound)
		}
		return fmt.Errorf("%s: exec query: %w", op, err)
	}
	return nil
}
