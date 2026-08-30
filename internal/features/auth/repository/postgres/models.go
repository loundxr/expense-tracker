package auth_postgres_repository

import "time"

type UserModel struct {
	ID           int
	Version      int
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	Role         string
}
