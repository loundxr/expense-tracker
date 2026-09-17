package transport

import (
	"context"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_http_middleware "github.com/loundxr/expense-tracker/internal/core/transport/http/middleware"
)

type UsersHTTPHandler struct {
	userService UsersService
	logger      *slog.Logger
}

type UsersService interface {
	GetUsers(ctx context.Context, limit *int, offset *int) ([]domain.User, error)
	GetUser(ctx context.Context, id int64) (domain.User, error)
	DeleteUser(ctx context.Context, id int64) error
	PatchUser(ctx context.Context, id int64, patch domain.UserPatch) (domain.User, error)
	UpdateUserRole(ctx context.Context, id int64, role string) (domain.User, error)
}

func NewUsersHTTPHandler(us UsersService, l *slog.Logger) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		userService: us,
		logger:      l,
	}
}

func (h *UsersHTTPHandler) RegisterRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/{id}", h.GetUser)
		r.Delete("/{id}", h.DeleteUser)
		r.Patch("/{id}", h.PatchUser)

		r.Group(func(r chi.Router) {
			r.Use(core_http_middleware.AdminOnly(h.logger))
			r.Get("/", h.GetUsers)
			r.Patch("/{id}/role", h.UpdateUserRole)
		})
	})
}
