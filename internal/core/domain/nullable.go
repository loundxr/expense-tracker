package domain

type Nullable[T any] struct {
	Val *T
	Set bool
}
