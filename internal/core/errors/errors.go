package core_errors

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidArgument  = errors.New("invalid argument")
	ErrAlreadyExists    = errors.New("entity already exists")
	ErrConflict         = errors.New("conflict")
	ErrNotFound         = errors.New("not found")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrCategoryNotFound = fmt.Errorf("%w: category", ErrNotFound)
	ErrAccountNotFound  = fmt.Errorf("%w: account", ErrNotFound)
)
