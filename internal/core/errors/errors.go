package core_errors

import "errors"

var (
	ErrInvalidArgument  = errors.New("invalid argument")
	ErrAlreadyExists    = errors.New("entity already exists")
	ErrConflict         = errors.New("conflict")
	ErrNotFound         = errors.New("not found")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrCategoryNotFound = errors.New("category not found")
	ErrAccountNotFound  = errors.New("account not found")
)
