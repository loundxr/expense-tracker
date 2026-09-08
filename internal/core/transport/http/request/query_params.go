package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
)

func GetIntQueryParams(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"param=%s by key=%s is not a valid integer: %v: %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return &val, nil
}

func ValidateLimitOffsetParams(limit *int, offset *int) error {
	if limit != nil && *limit <= 0 {
		return fmt.Errorf(
			"limit must be positive: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if limit != nil && *limit > 100 {
		*limit = 100
	}

	if offset != nil && *offset < 0 {
		return fmt.Errorf(
			"offset must be non-negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	return nil
}

func GetLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	const (
		limitQueryKey  = "limit"
		offsetQueryKey = "offset"
	)

	limit, err := GetIntQueryParams(r, limitQueryKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get %s query param: %w", limitQueryKey, err)
	}

	offset, err := GetIntQueryParams(r, offsetQueryKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get %s query param: %w", offsetQueryKey, err)
	}
	return limit, offset, nil
}
