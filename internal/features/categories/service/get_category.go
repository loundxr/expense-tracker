package service

import (
	"context"
	"fmt"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_ctx "github.com/loundxr/expense-tracker/internal/core/transport/http/context"
)

func (s *CategoriesService) GetCategory(ctx context.Context, id int64) (domain.Category, error) {
	const op = "categories.service.GetCategory"

	uid := core_ctx.GetUserID(ctx)
	role := core_ctx.GetUserRole(ctx)

	cat, err := s.categoriesRepository.GetCategory(ctx, id)
	if err != nil {
		return domain.Category{}, fmt.Errorf("%s: %w", op, err)
	}

	if role != domain.RoleAdmin {
		isSystem := cat.UserID == nil
		isOwner := !isSystem && *cat.UserID == uid

		if !isSystem && !isOwner {
			return domain.Category{}, fmt.Errorf("%s: %w", op, core_errors.ErrForbidden)
		}
	}

	return cat, nil
}
