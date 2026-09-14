package expenses_service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_ctx "github.com/loundxr/expense-tracker/internal/core/transport/http/context"
)

func (s *ExpensesService) GetExpenses(
	ctx context.Context,
	filter domain.ExpenseFilter,
) ([]domain.Expense, error) {
	const op = "expenses.service.GetExpenses"

	uid := core_ctx.GetUserID(ctx)
	role := core_ctx.GetUserRole(ctx)

	if role != domain.RoleAdmin {
		hasAccess, err := s.accChecker.HasAccess(ctx, uid, filter.AccountID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if !hasAccess {
			return nil, fmt.Errorf("%s: %w", op, core_errors.ErrForbidden)
		}
	}

	if filter.CategoryID != nil {
		hasAccess, err := s.catChecker.HasAccess(ctx, uid, *filter.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if !hasAccess {
			return nil, fmt.Errorf("%s: %w", op, core_errors.ErrCategoryNotFound)
		}
	}

	expenses, err := s.expensesRepository.GetExpenses(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	s.logger.Info(
		"user fetched expenses",
		slog.Int64("user_id", uid),
	)

	return expenses, nil
}
