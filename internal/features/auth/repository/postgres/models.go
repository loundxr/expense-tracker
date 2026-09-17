package repository

import "time"

type UserModel struct {
	ID           int64
	Version      int64
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	Role         string
}
