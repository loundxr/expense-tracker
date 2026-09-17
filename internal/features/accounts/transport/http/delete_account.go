package transport

import (
	"net/http"

	core_http_request "github.com/loundxr/expense-tracker/internal/core/transport/http/request"
	core_http_response "github.com/loundxr/expense-tracker/internal/core/transport/http/response"
)

func (h *AccountsHTTPHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rh := core_http_response.NewHTTPResponseHandler(w, h.logger)

	accountID, err := core_http_request.GetInt64IDPathValue(r, "id")
	if err != nil {
		rh.ErrorResponse(err, "invalid account 'id' in path")
		return
	}

	if err = h.accountsService.DeleteAccount(ctx, accountID); err != nil {
		rh.ErrorResponse(err, "failed to delete account")
		return
	}

	rh.NoContentResponse()
}
