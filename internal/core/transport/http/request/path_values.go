package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
)

func GetIntPathValue(r *http.Request, key string) (int, error) {
	pathVal := r.PathValue(key)
	if pathVal == "" {
		return 0, fmt.Errorf(
			"no key '%s' in path value: %w",
			key,
			core_errors.ErrInvalidArgument,
		)
	}

	val, err := strconv.Atoi(pathVal)
	if err != nil {
		return 0, fmt.Errorf(
			"path value=%s by key=%s is not a valid integer: %v: %w",
			pathVal,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}
	return val, nil
}
