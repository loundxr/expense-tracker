package repository

import (
	"context"
	"fmt"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

func (r *AnalyticsRepository) Trends(ctx context.Context, f domain.TrendsFilter) ([]domain.ExpenseTrendPoint, error) {
	const op = "analytics.repository.postgres.Trends"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	SELECT
		DATE_TRUNC($1, date) AS period,
		SUM(amount)::bigint AS total_amount,
		COUNT(*) AS expenses_count
	FROM expense_tracker.expenses
	WHERE account_id=$2
		AND ($3::timestamptz IS NULL OR date >= $3)
		AND ($4::timestamptz IS NULL OR date < $4)
	GROUP BY period
	ORDER BY period ASC;`

	rows, err := r.pool.Query(ctx, query, f.Interval, f.AccountID, f.From, f.To)
	if err != nil {
		return nil, fmt.Errorf("%s: execute query: %w", op, err)
	}
	defer rows.Close()

	points := make([]domain.ExpenseTrendPoint, 0)
	for rows.Next() {
		var p domain.ExpenseTrendPoint
		err := rows.Scan(
			&p.Date,
			&p.TotalAmount,
			&p.ExpensesCount,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan rows: %w", op, err)
		}
		points = append(points, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}

	return points, err
}
