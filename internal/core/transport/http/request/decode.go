package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
)

// TODO: rework to -> DecodeAndValidate
func DecodeJSONRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil { // check for .Decode(&dest)
		return fmt.Errorf("decode json: %v: %w", err, core_errors.ErrInvalidArgument)
	}
	return nil
}
