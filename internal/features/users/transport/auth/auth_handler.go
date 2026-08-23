package users_transport_http_auth

import (
	"context"
	"log/slog"

	"github.com/go-chi/chi"
	"github.com/loundxr/expense-tracker/internal/core/domain"
)

type AuthHandler struct {
	authService AuthService
	logger      *slog.Logger
}

type AuthService interface {
	SignUp(ctx context.Context, email, password string) (domain.User, error)
	SignIn(ctx context.Context, email, password string) (string, error)
}

func NewAuthHandler(service AuthService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		authService: service,
		logger:      logger,
	}
}

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", h.SignUp)
		//r.Post("/signin", h.SignIn)
	})
}
