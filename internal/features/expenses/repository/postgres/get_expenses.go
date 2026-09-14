package expenses_repository_postgres

import (
	"context"
	"fmt"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

func (r *ExpensesRepository) GetExpenses(ctx context.Context, filter domain.ExpenseFilter) ([]domain.Expense, error) {
	const op = "expenses.repository.postgres.GetExpenses"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	SELECT id, version, account_id, category_id,
		amount, currency, description, date, created_at
	FROM expense_tracker.expenses
	WHERE account_id=$1
		AND ($2::bigint IS NULL OR category_id=$2)
		AND ($3::timestamptz IS NULL OR date >= $3)
		AND ($4::timestamptz IS NULL OR date < $4)
	ORDER BY date, id DESC
	LIMIT $5 OFFSET $6;`

	rows, err := r.pool.Query(
		ctx, query,
		filter.AccountID,
		filter.CategoryID,
		filter.From, filter.To,
		filter.Limit, filter.Offset,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: execute query: %w", op, err)
	}
	defer rows.Close()

	models := make([]ExpenseModel, 0)
	for rows.Next() {
		var em ExpenseModel
		err := rows.Scan(
			&em.ID,
			&em.Version,
			&em.AccountID,
			&em.CategoryID,
			&em.Amount,
			&em.Currency,
			&em.Description,
			&em.Date,
			&em.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan returned row: %w", op, err)
		}
		models = append(models, em)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	expenses := expenseDomainsFromModels(models)
	return expenses, nil
}
