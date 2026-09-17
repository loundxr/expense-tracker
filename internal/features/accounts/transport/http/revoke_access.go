package transport

import (
	"net/http"

	core_http_request "github.com/loundxr/expense-tracker/internal/core/transport/http/request"
	core_http_response "github.com/loundxr/expense-tracker/internal/core/transport/http/response"
)

type RevokeAccessRequest struct {
	Email string `json:"email" validate:"email,required"`
}

func (h *AccountsHTTPHandler) RevokeAccess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rh := core_http_response.NewHTTPResponseHandler(w, h.logger)

	accountID, err := core_http_request.GetInt64IDPathValue(r, "id")
	if err != nil {
		rh.ErrorResponse(err, "invalid id value")
		return
	}

	var req RevokeAccessRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		rh.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	if err := h.accountsService.RevokeAccess(ctx, accountID, req.Email); err != nil {
		rh.ErrorResponse(err, "failed to revoke access to account")
		return
	}
	rh.NoContentResponse()
}
