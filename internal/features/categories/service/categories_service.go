package service

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
	GetCategoriesByID(ctx context.Context, uid int64) ([]domain.Category, error)
	GetCategory(ctx context.Context, id int64) (domain.Category, error)
	DeleteCategory(ctx context.Context, id int64) error
	PatchCategory(ctx context.Context, id int64, toPatch domain.Category) (domain.Category, error)
}

func NewCategoriesService(cr CategoriesRepository, l *slog.Logger) *CategoriesService {
	return &CategoriesService{
		categoriesRepository: cr,
		logger:               l,
	}
}
