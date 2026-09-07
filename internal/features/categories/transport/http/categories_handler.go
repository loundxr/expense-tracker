package categories_transport_http

import (
	"context"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/loundxr/expense-tracker/internal/core/domain"
)

type CategoriesHTTPHandler struct {
	categoriesService CategoriesService
	logger            *slog.Logger
}

type CategoriesService interface {
	CreateCategory(ctx context.Context, cat domain.Category) (domain.Category, error)
}

func NewCategoriesHTTPHandler(cs CategoriesService, l *slog.Logger) *CategoriesHTTPHandler {
	return &CategoriesHTTPHandler{
		categoriesService: cs,
		logger:            l,
	}
}

func (h *CategoriesHTTPHandler) RegisterRoutes(r chi.Router) {
	r.Route("/categories", func(r chi.Router) {
		r.Post("/", h.CreateCategory)
	})
}
