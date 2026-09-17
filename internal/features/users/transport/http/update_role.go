package transport

import (
	"net/http"

	core_http_request "github.com/loundxr/expense-tracker/internal/core/transport/http/request"
	core_http_response "github.com/loundxr/expense-tracker/internal/core/transport/http/response"
)

type UpdateRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=user admin"`
}

type UpdateRoleResponse UserDTOResponse

func (h *UsersHTTPHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rh := core_http_response.NewHTTPResponseHandler(w, h.logger)

	id, err := core_http_request.GetInt64IDPathValue(r, "id")
	if err != nil {
		rh.ErrorResponse(err, "failed to get 'id' path value")
		return
	}

	var req UpdateRoleRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		rh.ErrorResponse(err, "failed to decode and validate http request")
		return
	}

	user, err := h.userService.UpdateUserRole(ctx, id, req.Role)
	if err != nil {
		rh.ErrorResponse(err, "failed to update role")
		return
	}

	response := UpdateRoleResponse(userDTOFromDomain(user))
	rh.JSONResponse(response, http.StatusOK)
}
