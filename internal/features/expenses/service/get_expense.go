package expenses_service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_ctx "github.com/loundxr/expense-tracker/internal/core/transport/http/context"
)

func (s *ExpensesService) GetExpense(ctx context.Context, id int64) (domain.Expense, error) {
	const op = "expenses.service.GetExpense"

	uid := core_ctx.GetUserID(ctx)
	role := core_ctx.GetUserRole(ctx)

	exp, err := s.expensesRepository.GetExpense(ctx, id)
	if err != nil {
		return domain.Expense{}, fmt.Errorf("%s: %w", op, err)
	}

	if role != domain.RoleAdmin {
		hasAccess, err := s.accChecker.HasAccess(ctx, uid, exp.AccountID)
		if err != nil {
			return domain.Expense{}, fmt.Errorf("%s: %w", op, err)
		}
		if !hasAccess {
			return domain.Expense{}, fmt.Errorf("%s: %w", op, core_errors.ErrForbidden)
		}
	}

	s.logger.Info(
		"expense fetched",
		slog.Int64("user_id", uid),
		slog.Int64("expense_id", exp.ID),
	)
	return exp, nil
}
