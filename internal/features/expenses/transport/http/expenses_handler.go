package expenses_transport_http

import (
	"context"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/loundxr/expense-tracker/internal/core/domain"
)

type ExpensesHTTPHandler struct {
	expensesService ExpensesService
	logger          *slog.Logger
}

type ExpensesService interface {
	CreateExpense(ctx context.Context, expense domain.Expense) (domain.Expense, error)
}

func NewExpensesHTTPHandler(es ExpensesService, l *slog.Logger) *ExpensesHTTPHandler {
	return &ExpensesHTTPHandler{
		expensesService: es,
		logger:          l,
	}
}

func (h *ExpensesHTTPHandler) RegisterRoutes(r chi.Router) {
	r.Route("/expenses", func(r chi.Router) {
		r.Post("/", h.CreateExpense)
	})
}
