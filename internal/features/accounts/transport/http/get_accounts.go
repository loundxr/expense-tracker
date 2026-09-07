package accounts_transport_http

import (
	"net/http"

	core_http_request "github.com/loundxr/expense-tracker/internal/core/transport/http/request"
	core_http_response "github.com/loundxr/expense-tracker/internal/core/transport/http/response"
)

type GetAccountsResponse []AccountDTOResponse

func (h *AccountsHTTPHandler) GetAccounts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rh := core_http_response.NewHTTPResponseHandler(w, h.logger)

	limit, offset, err := core_http_request.GetLimitOffsetQueryParams(r)
	if err != nil {
		rh.ErrorResponse(err, "failed to get limit/offset query params")
		return
	}

	err = core_http_request.ValidateLimitOffsetParams(limit, offset)
	if err != nil {
		rh.ErrorResponse(err, "invalid limit/offset values")
		return
	}

	accounts, err := h.accountsService.GetAccounts(ctx, limit, offset)
	if err != nil {
		rh.ErrorResponse(err, "failed to get accounts")
		return
	}

	response := GetAccountsResponse(accountDTOsFromDomains(accounts))
	rh.JSONResponse(response, http.StatusOK)
}
