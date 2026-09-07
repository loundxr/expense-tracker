package categories_service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_ctx "github.com/loundxr/expense-tracker/internal/core/transport/http/context"
)

func (s CategoriesService) GetCategories(ctx context.Context) ([]domain.Category, error) {
	const op = "categories.service.GetCategories"

	uid := core_ctx.GetUserID(ctx)
	role := core_ctx.GetUserRole(ctx)

	var (
		categories []domain.Category
		err        error
	)

	var fetchType string
	if role == domain.RoleAdmin {
		categories, err = s.categoriesRepository.GetAllCategories(ctx)
		fetchType = "all"
	} else {
		categories, err = s.categoriesRepository.GetCategoriesByID(ctx, uid)
		fetchType = "by_id"
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	s.logger.Info(
		"user fetched categories",
		slog.Int("id", uid),
		slog.String("type", fetchType),
	)

	return categories, nil
}
