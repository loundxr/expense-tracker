package core_pgx_pool

import (
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

func (r pgxRow) Scan(dest ...any) error {
	if err := r.Row.Scan(dest...); err != nil {
		return mapErrors(err)
	}
	return nil
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
