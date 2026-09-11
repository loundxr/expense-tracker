package expenses_repository_postgres

import (
	"context"
	"fmt"

	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
)

func (r *ExpensesRepository) DeleteExpense(ctx context.Context, id int64) error {
	const op = "expenses.repository.postgres.DeleteExpense"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	DELETE FROM expense_tracker.expenses
	WHERE id=$1;`

	tag, err := r.pool.Exec(ctx, query, id)
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, core_errors.ErrNotFound)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
