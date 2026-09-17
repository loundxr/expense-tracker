package repository

import core_postgres_pool "github.com/loundxr/expense-tracker/internal/core/repository/postgres/pool"

type AccountsRepository struct {
	pool core_postgres_pool.Pool
}

func NewAccountsRepository(pool core_postgres_pool.Pool) *AccountsRepository {
	return &AccountsRepository{
		pool: pool,
	}
}
