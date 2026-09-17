package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_postgres_pool "github.com/loundxr/expense-tracker/internal/core/repository/postgres/pool"
)

func (r *AccountsRepository) CreateAccount(ctx context.Context, a domain.Account) (domain.Account, error) {
	const op = "accounts.repository.postgres.CreateAccount"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Account{}, fmt.Errorf("%s: begin tx: %w", op, err)
	}
	defer tx.Rollback(ctx)

	queryAccount := `
	INSERT INTO expense_tracker.accounts (name, user_id, created_at)
	VALUES ($1, $2, $3)
	RETURNING id, version, name, user_id, created_at;`

	row := r.pool.QueryRow(ctx, queryAccount, a.Name, a.UserID, a.CreatedAt)

	var am AccountModel
	err = row.Scan(
		&am.ID,
		&am.Version,
		&am.Name,
		&am.UserID,
		&am.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return domain.Account{}, fmt.Errorf(
				"%s: user with id=%d: %w",
				op,
				a.UserID,
				core_errors.ErrNotFound,
			)
		}
		return domain.Account{}, fmt.Errorf("%s: insert account: %w", op, err)
	}

	queryAccess := `
	INSERT INTO expense_tracker.account_users (account_id, user_id, created_at)
	VALUES($1, $2, $3);`

	_, err = r.pool.Exec(ctx, queryAccess, am.ID, am.UserID, am.CreatedAt)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return domain.Account{}, fmt.Errorf(
				"%s: user or account not found: %w",
				op,
				core_errors.ErrNotFound,
			)
		}
		return domain.Account{}, fmt.Errorf("%s: insert access: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Account{}, fmt.Errorf("%s: commit: %w", op, err)
	}
	accDomain := accountDomainFromModel(am)

	return accDomain, nil
}
