package domain

import (
	"fmt"
	"time"

	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
)

type User struct {
	ID           int64
	Version      int64
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
	id int64,
	version int64,
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

func (u *User) Validate() error {
	emailLen := len([]rune(u.Email))
	if emailLen < 3 || emailLen > 100 {
		return fmt.Errorf(
			"invalid 'email' length: %d: %w",
			emailLen,
			core_errors.ErrInvalidArgument,
		)
	}
	return nil
}

type UserPatch struct {
	Email    Nullable[string]
	Password Nullable[string]
}

func NewUserPatch(email Nullable[string], password Nullable[string]) UserPatch {
	return UserPatch{
		Email:    email,
		Password: password,
	}
}

func (p *UserPatch) Validate() error {
	if p.Email.Set && p.Email.Val == nil {
		return fmt.Errorf(
			"email cannot be changed to null: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.Password.Set && p.Password.Val == nil {
		return fmt.Errorf(
			"password cannot be changed to null: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	return nil
}

func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate patch: %w", err)
	}

	tmp := *u
	if patch.Email.Set {
		tmp.Email = *patch.Email.Val
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched user: %w", err)
	}

	*u = tmp
	return nil
}
