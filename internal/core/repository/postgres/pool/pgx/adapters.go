package core_pgx_pool

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	core_postgres_pool "github.com/loundxr/expense-tracker/internal/core/repository/postgres/pool"
)

type pgxRows struct {
	pgx.Rows
}

type pgxRow struct {
	pgx.Row
}

type pgconnCommandTag struct {
	pgconn.CommandTag
}

type pgxTx struct {
	pgx.Tx
}

func (t *pgxTx) QueryRow(ctx context.Context, sql string, args ...any) core_postgres_pool.Row {
	return pgxRow{
		Row: t.Tx.QueryRow(ctx, sql, args...),
	}
}

func (t *pgxTx) Exec(ctx context.Context, sql string, arguments ...any) (core_postgres_pool.CommandTag, error) {
	tag, err := t.Tx.Exec(ctx, sql, arguments...)
	if err != nil {
		return tag, mapErrors(err)
	}
	return pgconnCommandTag{tag}, nil
}

func (t *pgxTx) Commit(ctx context.Context) error {
	return t.Tx.Commit(ctx)
}

func (t *pgxTx) Rollback(ctx context.Context) error {
	return t.Tx.Rollback(ctx)
}

func (r pgxRow) Scan(dest ...any) error {
	if err := r.Row.Scan(dest...); err != nil {
		return mapErrors(err)
	}
	return nil
}

func (t pgconnCommandTag) RowsAffected() int64 {
	return t.CommandTag.RowsAffected()
}

func mapErrors(err error) error {
	const (
		pgxViolatesForeignKeyErrorCode       = "23503"
		pgxViolatesUniqueConstraintErrorCode = "23505"
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return core_postgres_pool.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgxViolatesForeignKeyErrorCode:
			return fmt.Errorf("%v: %w", err, core_postgres_pool.ErrViolatesForeignKey)
		case pgxViolatesUniqueConstraintErrorCode:
			return fmt.Errorf("%v: %w", err, core_postgres_pool.ErrViolatesUniqueConstraint)
		}
	}

	return fmt.Errorf("%v: %w", err, core_postgres_pool.ErrUnknown)
}
