package transport

import (
	"net/http"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_ctx "github.com/loundxr/expense-tracker/internal/core/transport/http/context"
	core_http_request "github.com/loundxr/expense-tracker/internal/core/transport/http/request"
	core_http_response "github.com/loundxr/expense-tracker/internal/core/transport/http/response"
)

type CreateAccountRequest struct {
	Name string `json:"name"`
}

type CreateAccountResponse AccountDTOResponse

func (h *AccountsHTTPHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rh := core_http_response.NewHTTPResponseHandler(w, h.logger)

	var req CreateAccountRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		rh.ErrorResponse(err, "invalid request")
		return
	}

	uninitAcc := domain.NewUninitializedAccount(req.Name, core_ctx.GetUserID(ctx))
	acc, err := h.accountsService.CreateAccount(ctx, uninitAcc)
	if err != nil {
		rh.ErrorResponse(err, "failed to create account")
		return
	}

	response := CreateAccountResponse(accountDTOFromDomain(acc))
	rh.JSONResponse(response, http.StatusCreated)
}
