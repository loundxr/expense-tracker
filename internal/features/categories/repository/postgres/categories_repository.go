package categories_repository_postgres

import (
	core_postgres_pool "github.com/loundxr/expense-tracker/internal/core/repository/postgres/pool"
)

type CategoriesRepository struct {
	pool core_postgres_pool.Pool
}

func NewCategoriesRepository(pool core_postgres_pool.Pool) *CategoriesRepository {
	return &CategoriesRepository{
		pool: pool,
	}
}
