package categories_service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

func (s *CategoriesService) CreateCategory(
	ctx context.Context,
	category domain.Category,
) (domain.Category, error) {
	const op = "categories.service.CreateCategory"

	cat, err := s.categoriesRepository.CreateCategory(ctx, category)
	if err != nil {
		return domain.Category{}, fmt.Errorf("%s: %w", op, err)
	}

	s.logger.Info(
		"user added personal category",
		slog.Int("user_id", *cat.UserID),
		slog.String("name", cat.Name),
	)
	return cat, nil
}
