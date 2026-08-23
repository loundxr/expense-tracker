package users_repository_postgres_auth

import core_postgres_pool "github.com/loundxr/expense-tracker/internal/core/repository/postgres/pool"

type UsersAuthRepository struct {
	pool core_postgres_pool.Pool
}

func NewUsersAuthRepository(pool core_postgres_pool.Pool) *UsersAuthRepository {
	return &UsersAuthRepository{
		pool: pool,
	}
}
