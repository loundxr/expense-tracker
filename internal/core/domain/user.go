package domain

import "time"

type User struct {
	ID           int
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	Version      int
}

func NewUser(
	id int,
	version int,
	email string,
	passwordHash string,
	createdAt time.Time,
) User {
	return User{
		ID:           id,
		Version:      version,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    createdAt,
	}
}

func NewUninitializedUser(email string, hash string) User {
	return User{
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    time.Now(),
		Version:      uninitialiedVersion,
		ID:           uninitializedID,
	}
}
