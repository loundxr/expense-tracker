package analytics_transport_http

import (
	"net/http"

	core_http_request "github.com/loundxr/expense-tracker/internal/core/transport/http/request"
	core_http_response "github.com/loundxr/expense-tracker/internal/core/transport/http/response"
)

type GetMembersContributionResponse MembersContributionDTOResponse

func (h *AnalyticsHTTPHandler) MembersContribution(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rh := core_http_response.NewHTTPResponseHandler(w, h.logger)

	accountID, err := core_http_request.GetInt64IDPathValue(r, "account_id")
	if err != nil {
		rh.ErrorResponse(err, "invalid 'account_id' path value")
		return
	}

	filter, err := getSummaryFilter(r, accountID)
	if err != nil {
		rh.ErrorResponse(err, "invalid from/to query params")
		return
	}

	members, err := h.analyticsService.MembersContribution(ctx, filter)
	if err != nil {
		rh.ErrorResponse(err, "failed to generate members contribution statistics")
		return
	}

	response := GetMembersContributionResponse(membersContributionDTOFromDomains(accountID, members))
	rh.JSONResponse(response, http.StatusOK)
}
