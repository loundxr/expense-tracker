package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_postgres_pool "github.com/loundxr/expense-tracker/internal/core/repository/postgres/pool"
)

func (r *ExpensesRepository) CreateExpense(ctx context.Context, expense domain.Expense) (domain.Expense, error) {
	const op = "expenses.repository.postgres.CreateExpense"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	INSERT INTO expense_tracker.expenses (
		account_id, 
		user_id, 
		category_id, 
		amount, 
		currency, 
		description, 
		date, 
		created_at
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id, version, account_id, user_id, category_id, amount, currency, description, date, created_at;`

	row := r.pool.QueryRow(ctx, query,
		expense.AccountID,
		expense.UserID,
		expense.CategoryID,
		expense.Amount,
		expense.Currency,
		expense.Description,
		expense.Date,
		expense.CreatedAt,
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
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return domain.Expense{}, fmt.Errorf("%s: %w", op, core_postgres_pool.MapFKErrors(err))
		}
		return domain.Expense{}, fmt.Errorf("%s: scan returned row: %w", op, err)
	}

	expDomain := expenseDomainFromModel(em)
	return expDomain, err
}
