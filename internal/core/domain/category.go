package domain

import "time"

type Category struct {
	ID        int
	Version   int
	Name      string
	UserID    *int
	CreatedAt time.Time
}

func NewCategory(
	id int,
	version int,
	name string,
	userID *int,
	createdAt time.Time,
) Category {
	return Category{
		ID:        id,
		Version:   version,
		Name:      name,
		UserID:    userID,
		CreatedAt: createdAt,
	}
}

func NewUninitializedCategory(name string, userID *int) Category {
	return NewCategory(
		uninitializedID,
		uninitialiedVersion,
		name,
		userID,
		time.Now(),
	)
}
