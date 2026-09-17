package analytics_transport_http

import (
	"fmt"
	"net/http"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_http_request "github.com/loundxr/expense-tracker/internal/core/transport/http/request"
	core_http_response "github.com/loundxr/expense-tracker/internal/core/transport/http/response"
)

func (h *AnalyticsHTTPHandler) Summary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rh := core_http_response.NewHTTPResponseHandler(w, h.logger)

	accountID, err := core_http_request.GetInt64IDPathValue(r, "account_id")
	if err != nil {
		rh.ErrorResponse(err, "invalid 'id' path value")
		return
	}

	filter, err := getSummaryFilter(r, accountID)
	if err != nil {
		rh.ErrorResponse(err, "invalid from/to query params")
		return
	}

	summary, err := h.analyticsService.Summary(ctx, filter)
	if err != nil {
		rh.ErrorResponse(err, "failed to generate financial summary")
		return
	}

	response := summaryDTOFromDomain(summary)
	rh.JSONResponse(response, http.StatusOK)
}

func getSummaryFilter(r *http.Request, id int64) (domain.SummaryFilter, error) {
	from, err := core_http_request.GetDateQueryParams(r, "from")
	if err != nil {
		return domain.SummaryFilter{}, err
	}

	to, err := core_http_request.GetDateQueryParams(r, "to")
	if err != nil {
		return domain.SummaryFilter{}, err
	}

	if from != nil && to != nil && from.After(*to) {
		return domain.SummaryFilter{},
			fmt.Errorf("parameter 'from' cannot be after 'to': %w", core_errors.ErrInvalidArgument)
	}

	return domain.NewSummaryFilter(id, from, to), nil
}
