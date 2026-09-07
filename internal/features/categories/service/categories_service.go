package categories_service

import (
	"context"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

type CategoriesService struct {
	categoriesRepository CategoriesRepository
	logger               *slog.Logger
}

type CategoriesRepository interface {
	CreateCategory(ctx context.Context, category domain.Category) (domain.Category, error)
	GetAllCategories(ctx context.Context) ([]domain.Category, error)
	GetCategoriesByID(ctx context.Context, uid int) ([]domain.Category, error)
}

func NewCategoriesService(cr CategoriesRepository, l *slog.Logger) *CategoriesService {
	return &CategoriesService{
		categoriesRepository: cr,
		logger:               l,
	}
}
