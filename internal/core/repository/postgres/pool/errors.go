package core_postgres_pool

import "fmt"

var (
	ErrNotFound           = fmt.Errorf("not found")
	ErrViolatesForeignKey = fmt.Errorf("violates foreign key")
	ErrUnknown            = fmt.Errorf("unknown error")
)
