package expenses_transport_http

import (
	"net/http"

	core_http_request "github.com/loundxr/expense-tracker/internal/core/transport/http/request"
	core_http_response "github.com/loundxr/expense-tracker/internal/core/transport/http/response"
)

func (h *ExpensesHTTPHandler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rh := core_http_response.NewHTTPResponseHandler(w, h.logger)

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil || id <= 0 {
		rh.ErrorResponse(err, "invalid 'id' path value")
		return
	}

	if err := h.expensesService.DeleteExpense(ctx, int64(id)); err != nil {
		rh.ErrorResponse(err, "failed to delete expense")
		return
	}

	rh.NoContentResponse()
}
