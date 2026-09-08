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
	GetCategories(ctx context.Context) ([]domain.Category, error)
	GetCategory(ctx context.Context, id int64) (domain.Category, error)
	DeleteCategory(ctx context.Context, id int64) error
	PatchCategory(ctx context.Context, id int64, patch domain.CategoryPatch) (domain.Category, error)
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
		r.Get("/", h.GetCategories)
		r.Get("/{id}", h.GetCategory)
		r.Delete("/{id}", h.DeleteCategory)
		r.Patch("/{id}", h.PatchCategory)
	})
}
