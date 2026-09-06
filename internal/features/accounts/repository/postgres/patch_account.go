package accounts_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_postgres_pool "github.com/loundxr/expense-tracker/internal/core/repository/postgres/pool"
)

func (r *AccountsRepository) PatchAccount(
	ctx context.Context,
	accountID int,
	acc domain.Account,
) (domain.Account, error) {
	const op = "accounts.repository.postgres.PatchAccount"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	UPDATE expense_tracker.accounts
	SET name=$1
	WHERE id=$2 AND version=$3
	RETURNING id, version, name, user_id, created_at`

	row := r.pool.QueryRow(ctx, query, acc.Name, accountID, acc.Version)

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
			return domain.Account{}, fmt.Errorf("%s: version mismatch: %w", op, core_errors.ErrConflict)
		}
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return domain.Account{}, fmt.Errorf("%s: owner user: %w", op, core_errors.ErrNotFound)
		}
		return domain.Account{}, fmt.Errorf("%s: %w", op, err)
	}

	accDomain := accountDomainFromModel(am)
	return accDomain, nil
}
