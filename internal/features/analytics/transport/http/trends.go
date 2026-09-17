package transport

import (
	"fmt"
	"net/http"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_http_request "github.com/loundxr/expense-tracker/internal/core/transport/http/request"
	core_http_response "github.com/loundxr/expense-tracker/internal/core/transport/http/response"
)

type GetTrendsResponse TrendsDTOResponse

func (h *AnalyticsHTTPHandler) Trends(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rh := core_http_response.NewHTTPResponseHandler(w, h.logger)

	accountID, err := core_http_request.GetInt64IDPathValue(r, "account_id")
	if err != nil {
		rh.ErrorResponse(err, "invalid 'account_id' path value")
		return
	}

	filter, err := getTrendsFilter(r, accountID)
	if err != nil {
		rh.ErrorResponse(err, "failed to get trends filter")
		return
	}

	points, err := h.analyticsService.Trends(ctx, filter)
	if err != nil {
		rh.ErrorResponse(err, "failed to fetch expense trends")
		return
	}

	response := GetTrendsResponse(trendsDTOResponseFromDomains(accountID, filter.Interval, points))

	rh.JSONResponse(response, http.StatusOK)
}

func getTrendsFilter(r *http.Request, accountID int64) (domain.TrendsFilter, error) {
	from, err := core_http_request.GetDateQueryParams(r, "from")
	if err != nil {
		return domain.TrendsFilter{}, err
	}

	to, err := core_http_request.GetDateQueryParams(r, "to")
	if err != nil {
		return domain.TrendsFilter{}, err
	}

	if from != nil && to != nil && from.After(*to) {
		return domain.TrendsFilter{},
			fmt.Errorf("parameter 'from' cannot be after 'to': %w", core_errors.ErrInvalidArgument)
	}

	interval := r.URL.Query().Get("interval")
	if interval == "" {
		interval = "day"
	}

	if interval != "day" && interval != "month" {
		return domain.TrendsFilter{},
			fmt.Errorf(
				"invalid interval '%s', expected 'day' or 'month': %w",
				interval,
				core_errors.ErrInvalidArgument,
			)
	}

	return domain.NewTrendsFilter(accountID, from, to, interval), nil
}
