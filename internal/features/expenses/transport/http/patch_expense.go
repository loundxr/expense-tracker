package expenses_transport_http

import (
	"net/http"
	"time"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_http_request "github.com/loundxr/expense-tracker/internal/core/transport/http/request"
	core_http_response "github.com/loundxr/expense-tracker/internal/core/transport/http/response"
	core_types "github.com/loundxr/expense-tracker/internal/core/types"
	"github.com/loundxr/expense-tracker/internal/utils/money"
)

type PatchExpenseRequest struct {
	Amount      core_types.Nullable[float64]   `json:"amount"`
	CategoryID  core_types.Nullable[int64]     `json:"category_id"`
	Description core_types.Nullable[string]    `json:"description"`
	Date        core_types.Nullable[time.Time] `json:"date"`
}

type PatchExpenseResponse ExpenseDTOResponse

func (h *ExpensesHTTPHandler) PatchExpense(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rh := core_http_response.NewHTTPResponseHandler(w, h.logger)

	id, err := core_http_request.GetInt64IDPathValue(r, "id")
	if err != nil {
		rh.ErrorResponse(err, "invalid 'id' path value format")
		return
	}

	var req PatchExpenseRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		rh.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	patch := expensePatchFromRequest(req)

	patchedExpense, err := h.expensesService.PatchExpense(ctx, id, patch)
	if err != nil {
		rh.ErrorResponse(err, "failed to patch expense")
		return
	}

	response := PatchExpenseResponse(expenseDTOFromDomain(patchedExpense))
	rh.JSONResponse(response, http.StatusOK)
}

func expensePatchFromRequest(req PatchExpenseRequest) domain.ExpensePatch {
	var centsAmount domain.Nullable[int64]
	if req.Amount.Set {
		inCents := money.ToCents(float64(*req.Amount.Val))
		centsAmount = domain.Nullable[int64]{
			Val: &inCents,
			Set: true,
		}

	}

	return domain.NewExpensePatch(
		centsAmount,
		req.CategoryID.ToDomain(),
		req.Description.ToDomain(),
		req.Date.ToDomain(),
	)
}
