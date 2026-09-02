package core_postgres_pool

import "errors"

var (
	ErrNotFound                 = errors.New("not found")
	ErrViolatesForeignKey       = errors.New("violates foreign key")
	ErrUnknown                  = errors.New("unknown error")
	ErrViolatesUniqueConstraint = errors.New("violates unique constraint")
)
