package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_ctx "github.com/loundxr/expense-tracker/internal/core/transport/http/context"
)

func (s *ExpensesService) CreateExpense(ctx context.Context, expense domain.Expense) (domain.Expense, error) {
	const op = "expenses.service.CreateExpense"

	uid := core_ctx.GetUserID(ctx)

	hasAccessToAcc, err := s.accChecker.HasAccess(ctx, uid, expense.AccountID)
	if err != nil {
		return domain.Expense{}, fmt.Errorf("%s: check account access: %w", op, err)
	}
	if !hasAccessToAcc {
		return domain.Expense{}, fmt.Errorf("%s: %w", op, core_errors.ErrForbidden)
	}

	hasAccessToCat, err := s.catChecker.HasAccess(ctx, uid, expense.CategoryID, expense.AccountID)
	if err != nil {
		return domain.Expense{}, fmt.Errorf("%s: check category access: %w", op, err)
	}
	if !hasAccessToCat {
		return domain.Expense{}, fmt.Errorf("%s: %w", op, core_errors.ErrCategoryNotFound)
	}

	expense.UserID = uid
	exp, err := s.expensesRepository.CreateExpense(ctx, expense)
	if err != nil {
		return domain.Expense{}, fmt.Errorf("%s: %w", op, err)
	}

	s.logger.Info(
		"expense created",
		slog.Int64("user_id", exp.UserID),
		slog.Int64("account_id", exp.AccountID),
		slog.Int64("category_id", exp.CategoryID),
		slog.Int64("amount", exp.Amount),
	)

	return exp, nil
}
