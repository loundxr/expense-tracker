package domain

import (
	"fmt"
	"time"

	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
)

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

func (c Category) Validate() error {
	nameLen := len([]rune(c.Name))
	if nameLen < 2 || nameLen > 100 {
		return fmt.Errorf(
			"invalid category name len: %d: %w",
			nameLen,
			core_errors.ErrInvalidArgument,
		)
	}
	return nil
}

type CategoryPatch struct {
	Name Nullable[string]
}

func NewCategoryPatch(name Nullable[string]) CategoryPatch {
	return CategoryPatch{
		Name: name,
	}
}

func (p CategoryPatch) Validate() error {
	if p.Name.Set && p.Name.Val == nil {
		return fmt.Errorf(
			"name of a category cannot be change to null: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	return nil
}

func (c *Category) ApplyPatch(patch CategoryPatch) error {
	if !patch.Name.Set {
		return fmt.Errorf("field 'name' is required: %w", core_errors.ErrInvalidArgument)
	}
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate patch: %w", err)
	}

	tmp := *c
	if patch.Name.Set {
		tmp.Name = *patch.Name.Val
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched category: %w", err)
	}
	*c = tmp
	return nil
}
