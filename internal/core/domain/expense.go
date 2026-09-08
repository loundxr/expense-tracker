package domain

import "time"

type Expense struct {
	ID          int64
	Version     int64
	AccountID   int64
	UserID      int64
	CategoryID  int64
	Amount      int64
	Currency    string
	Description string
	Date        time.Time
	CreatedAt   time.Time
}

func NewExpense(
	id int64,
	version int64,
	accountID int64,
	userID int64,
	categoryID int64,
	amount int64,
	currency string,
	description string,
	date time.Time,
	createdAt time.Time,
) Expense {
	return Expense{
		ID:          id,
		Version:     version,
		AccountID:   accountID,
		UserID:      userID,
		CategoryID:  categoryID,
		Amount:      amount,
		Currency:    currency,
		Description: description,
		Date:        date,
		CreatedAt:   createdAt,
	}
}

func NewUninitializedExpense(
	accountID, userID, categoryID int64,
	amount int64,
	description string,
	date time.Time,
) Expense {
	return NewExpense(
		uninitializedID,
		uninitialiedVersion,
		accountID,
		userID,
		categoryID,
		amount,
		"USD",
		description,
		date,
		time.Now(),
	)
}
