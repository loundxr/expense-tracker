package expenses_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_http_request "github.com/loundxr/expense-tracker/internal/core/transport/http/request"
	core_http_response "github.com/loundxr/expense-tracker/internal/core/transport/http/response"
	"github.com/loundxr/expense-tracker/internal/utils/money"
)

type CreateExpenseRequest struct {
	AccountID   int64     `json:"account_id"`
	CategoryID  int64     `json:"category_id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

type CreateExpenseResponse ExpenseDTOResponse

func (h *ExpensesHTTPHandler) CreateExpense(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rh := core_http_response.NewHTTPResponseHandler(w, h.logger)

	var req CreateExpenseRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		rh.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	amountInCents := money.ToCents(req.Amount)

	uninit := domain.NewUninitializedExpense(
		req.AccountID,
		req.CategoryID,
		amountInCents,
		req.Description,
		req.Date,
	)

	expense, err := h.expensesService.CreateExpense(ctx, uninit)
	if err != nil {
		rh.ErrorResponse(err, "failed to create expense")
		return
	}

	response := CreateExpenseResponse(expenseDTOFromDomain(expense))
	rh.JSONResponse(response, http.StatusCreated)
}

func (r *CreateExpenseRequest) Validate() error {
	if r.Date.IsZero() {
		r.Date = time.Now()
	}

	if r.Date.After(time.Now().Add(24 * time.Hour)) {
		return fmt.Errorf("you cannot set the date in the future: %w", core_errors.ErrInvalidArgument)
	}
	return nil
}
