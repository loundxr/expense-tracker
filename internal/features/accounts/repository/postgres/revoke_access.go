package accounts_repository_postgres

import (
	"context"
	"fmt"

	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
)

func (r *AccountsRepository) RevokeAccess(ctx context.Context, accountID int64, targetID int64) error {
	const op = "accounts.repository.postgres.RevokeAccess"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	DELETE FROM expense_tracker.account_users
	WHERE account_id=$1 AND user_id=$2;`

	tag, err := r.pool.Exec(ctx, query, accountID, targetID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf(
			"%s: user does not have access to this account yet: %w",
			op,
			core_errors.ErrNotFound,
		)
	}
	return nil
}
