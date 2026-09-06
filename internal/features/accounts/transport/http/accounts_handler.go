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
	GetAccount(ctx context.Context, id int) (domain.Account, error)
	DeleteAccount(ctx context.Context, accountID int) error
	PatchAccount(ctx context.Context, id int, patch domain.AccountPatch) (domain.Account, error)
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
		r.Get("/{id}", h.GetAccount)
		r.Delete("/{id}", h.DeleteAccount)
		r.Patch("/{id}", h.PatchAccount)
	})
}
