package analytics_transport_http

import (
	"net/http"

	core_http_request "github.com/loundxr/expense-tracker/internal/core/transport/http/request"
	core_http_response "github.com/loundxr/expense-tracker/internal/core/transport/http/response"
)

type CategoriesBreakdownResponse struct {
	AccountID  int64                  `json:"account_id"`
	Categories []CategoryBreakdownDTO `json:"categories_breakdown"`
}

func (h *AnalyticsHTTPHandler) CategoriesBreakdown(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rh := core_http_response.NewHTTPResponseHandler(w, h.logger)

	accountID, err := core_http_request.GetInt64IDPathValue(r, "account_id")
	if err != nil {
		rh.ErrorResponse(err, "invalid 'id' path value")
		return
	}

	filter, err := getFilter(r, accountID)
	if err != nil {
		rh.ErrorResponse(err, "invalid from/to query params")
		return
	}

	breakdown, err := h.analyticsService.CategoriesBreakdown(ctx, filter)
	if err != nil {
		rh.ErrorResponse(err, "failed to get categories breakdown")
		return
	}

	response := CategoriesBreakdownResponse{
		AccountID:  accountID,
		Categories: categoriesBreakdownDTOFromDomains(breakdown),
	}
	rh.JSONResponse(response, http.StatusOK)
}
