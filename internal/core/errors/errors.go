package core_errors

import "errors"

var (
	ErrInvalidArgument   = errors.New("invalid argument")
	ErrUserAlreadyExists = errors.New("user already exists")
)
