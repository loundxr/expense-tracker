package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_ctx "github.com/loundxr/expense-tracker/internal/core/transport/http/context"
)

func (s *CategoriesService) DeleteCategory(ctx context.Context, id int64) error {
	const op = "categories.service.DeleteCategory"

	uid := core_ctx.GetUserID(ctx)
	role := core_ctx.GetUserRole(ctx)

	category, err := s.categoriesRepository.GetCategory(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if role != domain.RoleAdmin {
		if category.UserID == nil {
			return fmt.Errorf("%s: cannot delete system category: %w", op, core_errors.ErrForbidden)
		}
		if uid != *category.UserID {
			return fmt.Errorf("%s: %w", op, core_errors.ErrForbidden)
		}
	}

	if err := s.categoriesRepository.DeleteCategory(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.logger.Info(
		"category deleted",
		slog.Int64("user_id", *category.UserID),
		slog.Int64("category_id", category.ID),
		slog.String("category", category.Name),
	)
	return nil
}
