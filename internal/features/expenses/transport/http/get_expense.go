package expenses_transport_http

import (
	"net/http"

	core_http_request "github.com/loundxr/expense-tracker/internal/core/transport/http/request"
	core_http_response "github.com/loundxr/expense-tracker/internal/core/transport/http/response"
)

type GetExpenseResponse ExpenseDTOResponse

func (h *ExpensesHTTPHandler) GetExpense(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rh := core_http_response.NewHTTPResponseHandler(w, h.logger)

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil || id <= 0 {
		rh.ErrorResponse(err, "invalid 'id' path value")
		return
	}

	expense, err := h.expensesService.GetExpense(ctx, int64(id))
	if err != nil {
		rh.ErrorResponse(err, "failed to get expense")
		return
	}

	response := GetExpenseResponse(expenseDTOFromDomain(expense))
	rh.JSONResponse(response, http.StatusOK)
}
