package accounts_transport_http

import (
	"net/http"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_http_request "github.com/loundxr/expense-tracker/internal/core/transport/http/request"
	core_http_response "github.com/loundxr/expense-tracker/internal/core/transport/http/response"
	core_types "github.com/loundxr/expense-tracker/internal/core/types"
)

type PatchAccountRequest struct {
	Name core_types.Nullable[string] `json:"name"`
}

type PatchAccountResponse AccountDTOResponse

func (h *AccountsHTTPHandler) PatchAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rh := core_http_response.NewHTTPResponseHandler(w, h.logger)

	accountID, err := core_http_request.GetInt64IDPathValue(r, "id")
	if err != nil {
		rh.ErrorResponse(err, "invalid id path value")
		return
	}

	var req PatchAccountRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		rh.ErrorResponse(err, "failed to decode and validate http request")
		return
	}

	patch := accountPatchFromRequest(req)
	accountDomain, err := h.accountsService.PatchAccount(ctx, accountID, patch)
	if err != nil {
		rh.ErrorResponse(err, "failed to patch account")
		return
	}

	response := PatchAccountResponse(accountDTOFromDomain(accountDomain))
	rh.JSONResponse(response, http.StatusOK)
}

func accountPatchFromRequest(req PatchAccountRequest) domain.AccountPatch {
	return domain.NewAccountPatch(req.Name.ToDomain())
}
