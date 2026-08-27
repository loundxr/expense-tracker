package utils

import "fmt"

func Must[T any](val T, err error) T {
	if err != nil {
		panic(fmt.Sprintf("critical error during initialization: %v", err))
	}
	return val
}
