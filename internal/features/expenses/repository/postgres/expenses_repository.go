package expenses_repository_postgres

import core_postgres_pool "github.com/loundxr/expense-tracker/internal/core/repository/postgres/pool"

type ExpensesRepository struct {
	pool core_postgres_pool.Pool
}

func NewExpensesRepository(pool core_postgres_pool.Pool) *ExpensesRepository {
	return &ExpensesRepository{
		pool: pool,
	}
}
