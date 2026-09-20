package core_postgres_pool

import (
	"context"
	"time"
)

type Pool interface {
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) Row
	Exec(ctx context.Context, sql string, arguments ...any) (CommandTag, error)
	Begin(ctx context.Context) (Transaction, error)
	WithinTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
	Close()
	OperationTimeout() time.Duration
}

type Transaction interface {
	QueryRow(ctx context.Context, sql string, args ...any) Row
	Exec(ctx context.Context, sql string, arguments ...any) (CommandTag, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type Transactor interface {
	WithinTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
}

type Rows interface {
	Close()
	Err() error
	Next() bool
	Scan(dest ...any) error
}

type Row interface {
	Scan(dest ...any) error
}

type CommandTag interface {
	RowsAffected() int64
}
