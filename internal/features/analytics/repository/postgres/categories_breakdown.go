package analytics_repository_postgres

import (
	"context"
	"fmt"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

func (r *AnalyticsRepository) CategoriesBreakdown(ctx context.Context, f domain.SummaryFilter) ([]domain.CategoryBreakdown, error) {
	const op = "analytics.repository.postgres.CategoriesBreakdown"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	SELECT
		COALESCE(c.id, 0) AS category_id,
		COALESCE(c.name, 'Unknown') AS category_name,
		SUM(e.amount)::bigint AS total_amount,
		COUNT(e.id) AS total_transactions,
		COALESCE(
			ROUND((SUM(e.amount)::numeric / NULLIF(SUM(SUM(e.amount)) OVER(), 0) * 100), 2),
			0
		) AS percentage
	FROM expense_tracker.expenses e
	LEFT JOIN expense_tracker.categories c ON e.category_id = c.id
	WHERE e.account_id=$1
		AND ($2::timestamptz IS NULL OR e.date >= $2)
		AND ($3::timestamptz IS NULL OR e.date < $3)
	GROUP BY c.id, c.name
	ORDER BY total_amount DESC;`

	rows, err := r.pool.Query(ctx, query, f.AccountID, f.From, f.To)
	if err != nil {
		return nil, fmt.Errorf("%s: execute query: %w", op, err)
	}
	defer rows.Close()

	res := make([]domain.CategoryBreakdown, 0)
	for rows.Next() {
		var cb domain.CategoryBreakdown
		err := rows.Scan(
			&cb.CategoryID,
			&cb.CategoryName,
			&cb.TotalAmount,
			&cb.TotalTransactions,
			&cb.Percentage,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan rows: %w", op, err)
		}
		res = append(res, cb)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}
	return res, nil
}
