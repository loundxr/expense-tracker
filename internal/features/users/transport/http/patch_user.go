package users_transport_http

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_http_request "github.com/loundxr/expense-tracker/internal/core/transport/http/request"
	core_http_response "github.com/loundxr/expense-tracker/internal/core/transport/http/response"
	core_types "github.com/loundxr/expense-tracker/internal/core/types"
)

type PatchUserRequest struct {
	Email    core_types.Nullable[string] `json:"email"`
	Password core_types.Nullable[string] `json:"password"`
}

type PatchUserResponse UserDTOResponse

func (h *UsersHTTPHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rh := core_http_response.NewHTTPResponseHandler(w, h.logger)

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		rh.ErrorResponse(err, "failed to get int path value 'id'")
		return
	}

	var req PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		rh.ErrorResponse(err, "failed to decode and validate http request")
		return
	}

	patch := userPatchFromRequest(req)
	userDomain, err := h.userService.PatchUser(ctx, id, patch)
	if err != nil {
		rh.ErrorResponse(err, "failed to patch user")
		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))
	rh.JSONResponse(response, http.StatusOK)
}

func (r *PatchUserRequest) Validate() error {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if r.Email.Set {
		if *r.Email.Val == "" {
			return fmt.Errorf("email cannot be empty:")
		}
		if !emailRegex.MatchString(*r.Email.Val) {
			return fmt.Errorf("invalid email format:")
		}
	}

	if r.Password.Set {
		if len([]rune(*r.Password.Val)) < 8 {
			return fmt.Errorf("new password must be longer than 8 symbols: %w", core_errors.ErrInvalidArgument)
		}
	}
	return nil
}

func userPatchFromRequest(req PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		req.Email.ToDomain(),
		req.Password.ToDomain(),
	)
}
