package users_transport_http

import (
	"context"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_http_middleware "github.com/loundxr/expense-tracker/internal/core/transport/http/middleware"
)

type UserHTTPHandler struct {
	userService UserService
	logger      *slog.Logger
}

type UserService interface {
	GetUsers(ctx context.Context, limit *int, offset *int) ([]domain.User, error)
}

func NewUserHTTPHandler(us UserService, l *slog.Logger) *UserHTTPHandler {
	return &UserHTTPHandler{
		userService: us,
		logger:      l,
	}
}

func (h *UserHTTPHandler) RegisterRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		// GetUser
		// DeleteUser
		// PatchUser

		r.Group(func(r chi.Router) {
			r.Use(core_http_middleware.AdminOnly(h.logger))
			r.Get("/", h.GetUsers)
		})
	})
}
