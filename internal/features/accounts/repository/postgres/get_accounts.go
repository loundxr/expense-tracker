package accounts_repository_postgres

import (
	"context"
	"fmt"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

func (r *AccountsRepository) GetAllAccounts(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.Account, error) {
	query := `
	SELECT id, version, name, user_id, created_at
	FROM expense_tracker.accounts
	ORDER BY id ASC
	LIMIT $1
	OFFSET $2;`

	return r.fetch(ctx, query, limit, offset)
}

func (r *AccountsRepository) GetAccountsByUserID(
	ctx context.Context,
	uid int,
	limit *int,
	offset *int,
) ([]domain.Account, error) {
	query := `
	SELECT a.id, a.version, a.name, a.user_id, a.created_at
	FROM expense_tracker.accounts a
	JOIN expense_tracker.account_users au ON a.id = au.account_id
	WHERE au.user_id = $1
	ORDER BY a.id ASC
	LIMIT $2
	OFFSET $3;`

	return r.fetch(ctx, query, uid, limit, offset)
}

func (r *AccountsRepository) fetch(
	ctx context.Context,
	query string,
	args ...any,
) ([]domain.Account, error) {
	const op = "accounts.repository.postgres.fetch"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: select accounts: %w", op, err)
	}
	defer rows.Close()

	var accountModels []AccountModel
	for rows.Next() {
		var am AccountModel
		err := rows.Scan(
			&am.ID,
			&am.Version,
			&am.Name,
			&am.UserID,
			&am.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: scan rows: %w", op, err)
		}
		accountModels = append(accountModels, am)
	}
	domains := accountDomainsFromModels(accountModels)
	return domains, nil
}
