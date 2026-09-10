package expenses_service

import (
	"context"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

type ExpensesService struct {
	expensesRepository ExpensesRepository
	accChecker         AccountAccessChecker
	catChecker         CategoryAccessChecker
	logger             *slog.Logger
}

func NewExpensesService(
	er ExpensesRepository,
	ac AccountAccessChecker,
	cc CategoryAccessChecker,
	l *slog.Logger,
) *ExpensesService {
	return &ExpensesService{
		expensesRepository: er,
		accChecker:         ac,
		catChecker:         cc,
		logger:             l,
	}
}

type ExpensesRepository interface {
	CreateExpense(ctx context.Context, expense domain.Expense) (domain.Expense, error)
}

type AccountAccessChecker interface {
	HasAccess(ctx context.Context, uid, accountID int64) (bool, error)
}

type CategoryAccessChecker interface {
	HasAccess(ctx context.Context, uid, categoryID int64) (bool, error)
}
