package accounts_repository_postgres

import (
	"context"
	"fmt"
)

func (r *AccountsRepository) DeleteAccount(ctx context.Context, accountID int) error {
	const op = "accounts.repository.postgres.DeleteAccount"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	DELETE FROM expense_tracker.accounts
	WHERE id=$1`

	_, err := r.pool.Exec(ctx, query, accountID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
