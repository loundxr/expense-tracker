package transport

import (
	"net/http"

	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_http_request "github.com/loundxr/expense-tracker/internal/core/transport/http/request"
	core_http_response "github.com/loundxr/expense-tracker/internal/core/transport/http/response"
)

func (h *UsersHTTPHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rh := core_http_response.NewHTTPResponseHandler(w, h.logger)

	uid, err := core_http_request.GetInt64IDPathValue(r, "id")
	if err != nil {
		rh.ErrorResponse(core_errors.ErrInvalidArgument, "invalid user id")
		return
	}

	if err := h.userService.DeleteUser(ctx, uid); err != nil {
		rh.ErrorResponse(err, "failed to delete user")
		return
	}
	rh.NoContentResponse()
}
