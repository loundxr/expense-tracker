package transport

import (
	"fmt"
	"net/http"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_http_request "github.com/loundxr/expense-tracker/internal/core/transport/http/request"
	core_http_response "github.com/loundxr/expense-tracker/internal/core/transport/http/response"
)

type GetExpensesResponse struct {
	Expenses []ExpenseDTOResponse `json:"expenses"`
}

func (h *ExpensesHTTPHandler) GetExpenses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rh := core_http_response.NewHTTPResponseHandler(w, h.logger)

	accountID, err := core_http_request.GetInt64IDPathValue(r, "account_id")
	if err != nil {
		rh.ErrorResponse(err, "invalid 'account_id' path value")
		return
	}

	filter, err := getFilter(r, accountID)
	if err != nil {
		rh.ErrorResponse(err, "invalid query params")
		return
	}

	expenses, err := h.expensesService.GetExpenses(ctx, filter)
	if err != nil {
		rh.ErrorResponse(err, "failed to get expenses list")
		return
	}

	response := NewGetExpensesResponse(expenseDTOsFromDomains(expenses))
	rh.JSONResponse(response, http.StatusOK)
}

func getFilter(r *http.Request, accountID int64) (domain.ExpenseFilter, error) {
	const (
		categoryIDKey = "category_id"
		fromKey       = "from"
		toKey         = "to"
	)
	var maxLimit = 100

	categoryID, err := core_http_request.GetInt64QueryParam(r, categoryIDKey)
	if err != nil {
		return domain.ExpenseFilter{}, err
	}

	from, err := core_http_request.GetDateQueryParams(r, fromKey)
	if err != nil {
		return domain.ExpenseFilter{}, err
	}

	to, err := core_http_request.GetDateQueryParams(r, toKey)
	if err != nil {
		return domain.ExpenseFilter{}, err
	}

	if from != nil && to != nil && from.After(*to) {
		return domain.ExpenseFilter{},
			fmt.Errorf("parameter 'from' cannot be after 'to': %w", core_errors.ErrInvalidArgument)
	}

	limit, offset, err := core_http_request.GetLimitOffsetQueryParams(r)
	if err != nil {
		return domain.ExpenseFilter{}, err
	}
	if limit != nil && *limit > 100 {
		limit = &maxLimit
	}

	return domain.NewExpenseFilter(
		accountID,
		categoryID,
		from,
		to,
		limit,
		offset,
	), nil
}

func NewGetExpensesResponse(expenses []ExpenseDTOResponse) GetExpensesResponse {
	if expenses == nil {
		expenses = make([]ExpenseDTOResponse, 0)
	}
	return GetExpensesResponse{
		Expenses: expenses,
	}
}
