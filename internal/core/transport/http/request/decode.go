package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
)

var reqValidator = validator.New()

type validatable interface {
	Validate() error
}

func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode json: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	if v, ok := dest.(validatable); ok {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("custom validation failed: %v: %w", err, core_errors.ErrInvalidArgument)
		}
	}

	if err := reqValidator.Struct(dest); err != nil {
		return fmt.Errorf("validate request struct: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	return nil
}
