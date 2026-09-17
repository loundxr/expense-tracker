package repository

import (
	"context"
	"fmt"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

func (r *AnalyticsRepository) MembersContribution(ctx context.Context, f domain.SummaryFilter) ([]domain.MemberContribution, error) {
	const op = "analytics.repository.postgres.MembersContribution"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	SELECT
		u.id AS user_id,
		u.email AS email,
		COALESCE(SUM(e.amount), 0)::bigint AS total_amount,
		COUNT(e.id) AS expenses_count,
		COALESCE(
			ROUND((SUM(e.amount)::numeric / NULLIF(SUM(SUM(e.amount)) OVER(), 0) * 100), 2),
			0
		)::float8 AS percentage
	FROM expense_tracker.account_users au
	JOIN expense_tracker.users u ON au.user_id = u.id
	LEFT JOIN expense_tracker.expenses e ON e.account_id = au.account_id
										AND e.user_id = au.user_id
										AND ($2::timestamptz IS NULL OR e.date >= $2)
										AND ($3::timestamptz IS NULL OR e.date < $3)
	WHERE au.account_id=$1
	GROUP BY u.id, u.email
	ORDER BY total_amount DESC, u.id ASC;`

	rows, err := r.pool.Query(ctx, query, f.AccountID, f.From, f.To)
	if err != nil {
		return nil, fmt.Errorf("%s: execute query: %w", op, err)
	}
	defer rows.Close()

	members := make([]domain.MemberContribution, 0)

	for rows.Next() {
		var m domain.MemberContribution
		err := rows.Scan(
			&m.UserID,
			&m.Email,
			&m.TotalAmount,
			&m.ExpensesCount,
			&m.Percentage,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan rows: %w", op, err)
		}
		members = append(members, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}
	return members, nil
}
