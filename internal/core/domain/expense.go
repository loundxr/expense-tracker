package domain

import (
	"fmt"
	"time"

	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
)

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
	accountID, categoryID int64,
	amount int64,
	description string,
	date time.Time,
) Expense {
	return NewExpense(
		uninitializedID,
		uninitialiedVersion,
		accountID,
		uninitializedID,
		categoryID,
		amount,
		"USD",
		description,
		date,
		time.Now(),
	)
}

func (e Expense) Validate() error {
	descLen := len([]rune(e.Description))
	if descLen < 3 {
		return fmt.Errorf(
			"invalid description length: %d: %w",
			descLen,
			core_errors.ErrInvalidArgument,
		)
	}

	if e.Date.After(e.CreatedAt.Add(24 * time.Hour)) {
		return fmt.Errorf(
			"you cannot set 'date' in the future: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if e.Amount <= 0 {
		return fmt.Errorf(
			"amount must be a positive number: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

type ExpensePatch struct {
	Amount      Nullable[int64]
	CategoryID  Nullable[int64]
	Description Nullable[string]
	Date        Nullable[time.Time]
}

func NewExpensePatch(
	amount Nullable[int64],
	categoryID Nullable[int64],
	description Nullable[string],
	date Nullable[time.Time],
) ExpensePatch {
	return ExpensePatch{
		Amount:      amount,
		CategoryID:  categoryID,
		Description: description,
		Date:        date,
	}
}

func (p ExpensePatch) Validate() error {
	if p.Amount.Set && p.Amount.Val == nil {
		return fmt.Errorf("amount cannot be set to null: %w", core_errors.ErrInvalidArgument)
	}

	if p.CategoryID.Set && p.CategoryID.Val == nil {
		return fmt.Errorf("category_id cannot be set to null: %w", core_errors.ErrInvalidArgument)
	}

	if p.Date.Set && p.Date.Val == nil {
		return fmt.Errorf("date cannot be set to null: %w", core_errors.ErrInvalidArgument)
	}
	return nil
}

func (e *Expense) ApplyPatch(p ExpensePatch) error {
	if err := p.Validate(); err != nil {
		return fmt.Errorf("validate patch: %w", err)
	}

	tmp := *e
	if p.Amount.Set && p.Amount.Val != nil {
		tmp.Amount = *p.Amount.Val
	}
	if p.CategoryID.Set && p.CategoryID.Val != nil {
		tmp.CategoryID = *p.CategoryID.Val
	}
	if p.Description.Set && p.Description.Val != nil {
		if p.Description.Val == nil {
			tmp.Description = ""
		} else {
			tmp.Description = *p.Description.Val
		}
	}
	if p.Date.Set && p.Date.Val != nil {
		tmp.Date = *p.Date.Val
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched entity: %w", err)
	}
	*e = tmp
	return nil
}

type ExpenseFilter struct {
	AccountID  int64
	CategoryID *int64
	From       *time.Time
	To         *time.Time
	Limit      *int
	Offset     *int
}

func NewExpenseFilter(accID int64, catID *int64, from, to *time.Time, limit, offset *int) ExpenseFilter {
	return ExpenseFilter{
		AccountID:  accID,
		CategoryID: catID,
		From:       from,
		To:         to,
		Limit:      limit,
		Offset:     offset,
	}
}
