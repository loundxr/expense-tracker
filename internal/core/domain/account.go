package domain

import "time"

type Account struct {
	ID        int
	Version   int
	Name      string
	UserID    int
	CreatedAt time.Time
}

func NewAccount(id int, version int, name string, userID int, createdAt time.Time) Account {
	return Account{
		ID:        id,
		Version:   version,
		Name:      name,
		UserID:    userID,
		CreatedAt: createdAt,
	}
}

func NewUninitializedAccount(name string, ownerID int) Account {
	return NewAccount(
		uninitializedID,
		uninitialiedVersion,
		name,
		ownerID,
		time.Now(),
	)
}
