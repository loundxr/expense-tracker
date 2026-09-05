package accounts_transport_http

import (
	"context"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/loundxr/expense-tracker/internal/core/domain"
)

type AccountsHTTPHandler struct {
	accountsService AccountsService
	logger          *slog.Logger
}

type AccountsService interface {
	CreateAccount(ctx context.Context, account domain.Account) (domain.Account, error)
	GetAccounts(ctx context.Context, limit *int, offset *int) ([]domain.Account, error)
}

func NewAccountsHTTPHandler(as AccountsService, l *slog.Logger) *AccountsHTTPHandler {
	return &AccountsHTTPHandler{
		accountsService: as,
		logger:          l,
	}
}

func (h *AccountsHTTPHandler) RegisterRoutes(r chi.Router) {
	r.Route("/accounts", func(r chi.Router) {
		r.Post("/", h.CreateAccount)
		r.Get("/", h.GetAccounts)
	})
}
