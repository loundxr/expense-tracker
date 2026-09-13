package expenses_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_postgres_pool "github.com/loundxr/expense-tracker/internal/core/repository/postgres/pool"
)

func (r *ExpensesRepository) PatchExpense(ctx context.Context, id int64, toPatch domain.Expense) (domain.Expense, error) {
	const op = "expenses.repository.postgres.PatchExpense"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	UPDATE expense_tracker.expenses
	SET
		amount=$1,
		category_id=$2,
		description=$3,
		date=$4,
		version=version+1
	WHERE id=$5 AND version=$6
	RETURNING
		id, version, account_id, user_id, category_id, amount, currency, description, date, created_at;`

	row := r.pool.QueryRow(
		ctx, query,
		toPatch.Amount,
		toPatch.CategoryID,
		toPatch.Description,
		toPatch.Date,
		id,
		toPatch.Version,
	)

	var em ExpenseModel
	err := row.Scan(
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
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNotFound) {
			return domain.Expense{}, fmt.Errorf("%s: version mismatch: %w", op, core_errors.ErrConflict)
		}
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return domain.Expense{}, fmt.Errorf("%s: %w", op, core_postgres_pool.MapFKErrors(err))
		}
		return domain.Expense{}, fmt.Errorf("%s: scan row: %w", op, err)
	}

	expenseDomain := expenseDomainFromModel(em)
	return expenseDomain, nil
}
