package domain

import "time"

type User struct {
	ID           int
	Version      int
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	Role         string
}

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

func NewUser(
	id int,
	version int,
	email string,
	passwordHash string,
	createdAt time.Time,
	role string,
) User {
	return User{
		ID:           id,
		Version:      version,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    createdAt,
		Role:         role,
	}
}

func NewUninitializedUser(email string, hash string) User {
	return User{
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    time.Now(),
		Version:      uninitialiedVersion,
		ID:           uninitializedID,
		Role:         uninitializedRole,
	}
}
