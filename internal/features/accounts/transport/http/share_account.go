package accounts_transport_http

import (
	"net/http"

	core_http_request "github.com/loundxr/expense-tracker/internal/core/transport/http/request"
	core_http_response "github.com/loundxr/expense-tracker/internal/core/transport/http/response"
)

type ShareAccountRequest struct {
	Email string `json:"email" validate:"email,required"`
}

func (h *AccountsHTTPHandler) ShareAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rh := core_http_response.NewHTTPResponseHandler(w, h.logger)

	uid, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil || uid < 0 {
		rh.ErrorResponse(err, "invalid id value")
		return
	}

	var req ShareAccountRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		rh.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	if err := h.accountsService.ShareAccount(ctx, uid, req.Email); err != nil {
		rh.ErrorResponse(err, "failed to share account")
		return
	}

	rh.NoContentResponse()
}
