package expenses_service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_ctx "github.com/loundxr/expense-tracker/internal/core/transport/http/context"
)

func (s *ExpensesService) DeleteExpense(ctx context.Context, id int64) error {
	const op = "expenses.service.DeleteExpense"

	uid := core_ctx.GetUserID(ctx)
	role := core_ctx.GetUserRole(ctx)

	expense, err := s.expensesRepository.GetExpense(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if role != domain.RoleAdmin && uid != expense.UserID {
		isOwner, err := s.accChecker.IsOwner(ctx, uid, expense.AccountID)
		if err != nil {
			return fmt.Errorf("%s: check account owner: %w", op, err)
		}
		if !isOwner {
			return fmt.Errorf("%s: %w", op, core_errors.ErrForbidden)
		}
	}

	if err := s.expensesRepository.DeleteExpense(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.logger.Info(
		"expense deleted",
		slog.Int64("user_id", uid),
		slog.Int64("account_id", expense.AccountID),
		slog.Int64("expense_id", expense.ID),
	)
	return nil
}
