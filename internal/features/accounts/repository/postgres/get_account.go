package accounts_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_postgres_pool "github.com/loundxr/expense-tracker/internal/core/repository/postgres/pool"
)

func (r *AccountsRepository) GetAccountByID(ctx context.Context, id int64) (domain.Account, error) {
	const op = "accounts.repository.postgres.GetAccountByID"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	SELECT id, version, name, user_id, created_at
	FROM expense_tracker.accounts
	WHERE id = $1;`

	row := r.pool.QueryRow(ctx, query, id)

	var am AccountModel

	err := row.Scan(
		&am.ID,
		&am.Version,
		&am.Name,
		&am.UserID,
		&am.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNotFound) {
			return domain.Account{}, fmt.Errorf("%s: %w", op, core_errors.ErrNotFound)
		}
		return domain.Account{}, fmt.Errorf("%s: scan row: %w", op, err)
	}

	accDomain := accountDomainFromModel(am)
	return accDomain, nil
}

func (r *AccountsRepository) HasAccess(ctx context.Context, uid, accountID int64) (bool, error) {
	const op = "accounts.repository.postgres.HasAccess"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	SELECT EXISTS(
		SELECT 1 FROM expense_tracker.account_users
		WHERE user_id=$1 AND account_id=$2);`

	var exists bool
	if err := r.pool.QueryRow(ctx, query, uid, accountID).Scan(&exists); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return exists, nil
}

func (r *AccountsRepository) IsOwner(ctx context.Context, uid, accountID int64) (bool, error) {
	const op = "accounts.repository.postgres.IsOwner"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	SELECT EXISTS(
		SELECT 1 FROM expense_tracker.accounts
		WHERE user_id=$1 AND id=$2)`

	var isOwner bool
	if err := r.pool.QueryRow(ctx, query, uid, accountID).Scan(&isOwner); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return isOwner, nil
}

func (r *AccountsRepository) Exists(ctx context.Context, id int64) error {
	const op = "accounts.repository.postgres.Exists"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	SELECT 1 FROM expense_tracker.accounts
	WHERE id=$1;`

	var tmp int
	if err := r.pool.QueryRow(ctx, query, id).Scan(&tmp); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNotFound) {
			return fmt.Errorf("%s: %w", op, core_errors.ErrNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
