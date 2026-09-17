package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_ctx "github.com/loundxr/expense-tracker/internal/core/transport/http/context"
)

func (s *CategoriesService) PatchCategory(
	ctx context.Context,
	id int64,
	patch domain.CategoryPatch,
) (domain.Category, error) {
	const op = "categories.service.PatchCategory"

	uid := core_ctx.GetUserID(ctx)
	role := core_ctx.GetUserRole(ctx)

	category, err := s.categoriesRepository.GetCategory(ctx, id)
	if err != nil {
		return domain.Category{}, fmt.Errorf("%s: %w", op, err)
	}

	if role != domain.RoleAdmin {
		if category.UserID == nil {
			return domain.Category{}, fmt.Errorf(
				"%s: cannot patch system categories: %w",
				op,
				core_errors.ErrForbidden,
			)
		}
		if *category.UserID != uid {
			return domain.Category{}, fmt.Errorf(
				"%s: %w",
				op,
				core_errors.ErrForbidden,
			)
		}
	}

	if err := category.ApplyPatch(patch); err != nil {
		return domain.Category{}, fmt.Errorf("%s: %w", op, err)
	}

	patched, err := s.categoriesRepository.PatchCategory(ctx, id, category)
	if err != nil {
		return domain.Category{}, fmt.Errorf("%s: %w", op, err)
	}

	s.logger.Info(
		"category patched",
		slog.Int64("user_id", uid),
		slog.Int64("category_id", patched.ID),
		slog.String("new_name", patched.Name),
	)

	return patched, nil
}
