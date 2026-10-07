package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_ctx "github.com/loundxr/expense-tracker/internal/core/transport/http/context"
)

func (s *ExpensesService) PatchExpense(ctx context.Context, id int64, p domain.ExpensePatch) (domain.Expense, error) {
	const op = "expenses.service.PatchExpense"

	expense, err := s.expensesRepository.GetExpense(ctx, id)
	if err != nil {
		return domain.Expense{}, fmt.Errorf("%s: %w", op, err)
	}

	uid := core_ctx.GetUserID(ctx)
	role := core_ctx.GetUserRole(ctx)

	if role != domain.RoleAdmin && uid != expense.UserID {
		isOwner, err := s.accChecker.IsOwner(ctx, uid, expense.AccountID)
		if err != nil || !isOwner {
			return domain.Expense{}, fmt.Errorf("%s: %w", op, core_errors.ErrForbidden)
		}
	}

	if p.CategoryID.Set && p.CategoryID.Val != nil {
		hasAccess, err := s.catChecker.HasAccess(ctx, uid, *p.CategoryID.Val, expense.AccountID)
		if err != nil || !hasAccess {
			return domain.Expense{}, fmt.Errorf("%s: %w", op, core_errors.ErrCategoryNotFound)
		}
	}

	if err := expense.ApplyPatch(p); err != nil {
		return domain.Expense{}, fmt.Errorf("%s: %w", op, err)
	}

	patched, err := s.expensesRepository.PatchExpense(ctx, id, expense)
	if err != nil {
		return domain.Expense{}, fmt.Errorf("%s: %w", op, err)
	}

	s.logger.Info(
		"expense patched",
		slog.Int64("expense_id", id),
		slog.Int64("user_id", uid),
		slog.Int64("account_id", patched.AccountID),
	)

	return patched, nil
}
