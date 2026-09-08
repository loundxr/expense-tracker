package domain

import (
	"fmt"
	"time"

	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
)

type Account struct {
	ID        int64
	Version   int64
	Name      string
	UserID    int64
	CreatedAt time.Time
}

func NewAccount(
	id int64,
	version int64,
	name string,
	userID int64,
	createdAt time.Time,
) Account {
	return Account{
		ID:        id,
		Version:   version,
		Name:      name,
		UserID:    userID,
		CreatedAt: createdAt,
	}
}

func NewUninitializedAccount(name string, ownerID int64) Account {
	return NewAccount(
		uninitializedID,
		uninitialiedVersion,
		name,
		ownerID,
		time.Now(),
	)
}

func (a *Account) Validate() error {
	nameLen := len([]rune(a.Name))
	if nameLen < 2 || nameLen > 100 {
		return fmt.Errorf(
			"inavlid 'name' length: %d: %w",
			nameLen,
			core_errors.ErrInvalidArgument,
		)
	}
	return nil
}

type AccountPatch struct {
	Name Nullable[string]
}

func NewAccountPatch(name Nullable[string]) AccountPatch {
	return AccountPatch{
		Name: name,
	}
}

func (p *AccountPatch) Validate() error {
	if p.Name.Set && p.Name.Val == nil {
		return fmt.Errorf("name of account cannot be changed to null: %w", core_errors.ErrInvalidArgument)
	}
	return nil
}

func (a *Account) ApplyPatch(patch AccountPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate patch: %w", err)
	}

	tmp := *a
	if patch.Name.Set {
		tmp.Name = *patch.Name.Val
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched account: %w", err)
	}
	*a = tmp
	return nil
}
