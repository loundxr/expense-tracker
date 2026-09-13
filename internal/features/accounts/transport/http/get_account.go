package accounts_transport_http

import (
	"net/http"

	core_http_request "github.com/loundxr/expense-tracker/internal/core/transport/http/request"
	core_http_response "github.com/loundxr/expense-tracker/internal/core/transport/http/response"
)

type GetAccountResponse AccountDTOResponse

func (h *AccountsHTTPHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rh := core_http_response.NewHTTPResponseHandler(w, h.logger)

	id, err := core_http_request.GetInt64IDPathValue(r, "id")
	if err != nil {
		rh.ErrorResponse(err, "invalid account 'id' value")
		return
	}

	acc, err := h.accountsService.GetAccount(ctx, id)
	if err != nil {
		rh.ErrorResponse(err, "failed to get account")
		return
	}

	response := GetAccountResponse(accountDTOFromDomain(acc))
	rh.JSONResponse(response, http.StatusOK)
}
