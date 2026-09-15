package analytics_repository_postgres

import (
	"context"
	"fmt"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

func (r *AnalyticsRepository) Summary(
	ctx context.Context,
	f domain.SummaryFilter,
) (domain.Summary, error) {
	const op = "analytics.repository.postgres.Summary"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	queryTotals := `
	SELECT
		COALESCE(SUM(amount), 0) AS total_amount,
		COUNT(*) AS total_transactions,
		COALESCE(AVG(amount), 0)::bigint AS avg_amount
	FROM expense_tracker.expenses
	WHERE account_id=$1
		AND ($2::timestamptz IS NULL OR date >= $2)
		AND ($3::timestamptz IS NULL OR date < $3);`

	s := domain.NewUninitializedSummary(f.AccountID)

	err := r.pool.QueryRow(ctx, queryTotals, f.AccountID, f.From, f.To).Scan(
		&s.TotalAmount,
		&s.TotalTransactions,
		&s.AverageAmount,
	)
	if err != nil {
		return domain.Summary{}, fmt.Errorf("%s: %w", op, err)
	}

	if s.TotalTransactions == 0 {
		return s, nil
	}

	queryMax := `
	SELECT id, amount, description, date
	FROM expense_tracker.expenses
	WHERE account_id=$1
		AND ($2::timestamptz IS NULL OR date >= $2)
		AND ($3::timestamptz IS NULL OR date < $3)
	ORDER BY amount DESC, date DESC
	LIMIT 1;`

	var maxExp domain.MaxExpenseSummary
	err = r.pool.QueryRow(ctx, queryMax, f.AccountID, f.From, f.To).Scan(
		&maxExp.ID,
		&maxExp.Amount,
		&maxExp.Description,
		&maxExp.Date,
	)

	if err == nil {
		s.MaxExpense = &maxExp
	}

	queryTopCat := `
	SELECT c.id, c.name, SUM(e.amount)::biging AS cat_total
	FROM expense_tracker.categories c
	JOIN expense_tracker.expenses e ON c.id = e.category_id
	WHERE e.account_id=$1
		AND ($2::timestamptz IS NULL OR date >= $2)
		AND ($3::timestamptz IS NULL OR date < $3);
	GROUP BY c.id, c.name
	ORDER BY cat_total DESC
	LIMIT 1;`

	var topCat domain.TopCategorySummary
	err = r.pool.QueryRow(ctx, queryTopCat, f.AccountID, f.From, f.To).Scan(
		&topCat.CategoryID,
		&topCat.CategoryName,
		&topCat.TotalAmount,
	)
	if err == nil {
		s.TopCategory = &topCat
	}

	return s, nil
}
