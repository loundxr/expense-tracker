package expenses_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_postgres_pool "github.com/loundxr/expense-tracker/internal/core/repository/postgres/pool"
)

func (r *ExpensesRepository) GetExpense(ctx context.Context, id int64) (domain.Expense, error) {
	const op = "expenses.repository.postgres.GetExpense"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	SELECT id, version, account_id, user_id, category_id,
		amount, currency, description, date, created_at
	FROM expense_tracker.expenses
	WHERE id=$1;`

	var em ExpenseModel
	if err := r.pool.QueryRow(ctx, query, id).Scan(
		&em.ID,
		&em.Version,
		&em.AccountID,
		&em.UserID,
		&em.CategoryID,
		&em.Amount,
		&em.Currency,
		&em.Description,
		&em.Date,
		&em.CreatedAt,
	); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNotFound) {
			return domain.Expense{}, fmt.Errorf("%s: %w", op, core_errors.ErrNotFound)
		}
		return domain.Expense{}, fmt.Errorf("%s: scan returned row: %w", op, err)
	}

	expenseDomain := expenseDomainFromModel(em)
	return expenseDomain, nil
}
