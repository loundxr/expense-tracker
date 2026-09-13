package users_transport_http

import (
	"net/http"

	core_http_request "github.com/loundxr/expense-tracker/internal/core/transport/http/request"
	core_http_response "github.com/loundxr/expense-tracker/internal/core/transport/http/response"
)

type GetUserResponse UserDTOResponse

func (h *UsersHTTPHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rh := core_http_response.NewHTTPResponseHandler(w, h.logger)

	uid, err := core_http_request.GetInt64IDPathValue(r, "id")
	if err != nil {
		rh.ErrorResponse(err, "invalid user id")
		return
	}

	user, err := h.userService.GetUser(ctx, uid)
	if err != nil {
		rh.ErrorResponse(err, "failed to get user")
		return
	}

	dto := GetUserResponse(userDTOFromDomain(user))
	rh.JSONResponse(dto, http.StatusOK)
}
