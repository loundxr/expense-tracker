package core_postgres_pool

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
)

var (
	ErrNotFound                 = errors.New("not found")
	ErrViolatesForeignKey       = errors.New("violates foreign key")
	ErrUnknown                  = errors.New("unknown error")
	ErrViolatesUniqueConstraint = errors.New("violates unique constraint")
)

const (
	fkExpensesAccount  = "expenses_account_id_fkey"
	fkExpensesCategory = "expenses_category_id_fkey"
)

func MapFKErrors(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.ConstraintName {
		case fkExpensesAccount:
			return core_errors.ErrAccountNotFound
		case fkExpensesCategory:
			return core_errors.ErrCategoryNotFound
		}
	}
	return err
}
