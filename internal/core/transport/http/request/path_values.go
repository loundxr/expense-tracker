package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
)

func GetInt64IDPathValue(r *http.Request, key string) (int64, error) {
	pathVal := r.PathValue(key)
	if pathVal == "" {
		return 0, fmt.Errorf(
			"no key '%s' in path value: %w",
			key,
			core_errors.ErrInvalidArgument,
		)
	}

	val, err := strconv.ParseInt(pathVal, 10, 64)
	if err != nil {
		return 0, fmt.Errorf(
			"path value=%s by key=%s is not a valid integer: %v: %w",
			pathVal,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	if val <= 0 {
		return 0, fmt.Errorf(
			"path param '%s' must be a positive integer, got %d: %w",
			key, val, core_errors.ErrInvalidArgument,
		)
	}
	return val, nil
}
