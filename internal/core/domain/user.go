package domain

import "time"

type User struct {
	ID           int
	Version      int
	Email        string
	PasswordHash string
	CreatedAt    time.Time
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
