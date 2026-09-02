package core_types

import (
	"encoding/json"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

type Nullable[T any] struct {
	domain.Nullable[T]
}

func (n *Nullable[T]) ToDomain() domain.Nullable[T] {
	return domain.Nullable[T]{
		Val: n.Val,
		Set: n.Set,
	}
}

func (n *Nullable[T]) UnmarshalJSON(b []byte) error {
	n.Set = true

	var val T
	if err := json.Unmarshal(b, &val); err != nil {
		return err
	}
	n.Val = &val
	return nil
}
