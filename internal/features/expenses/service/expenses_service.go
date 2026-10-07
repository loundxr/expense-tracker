package service

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
	GetExpenses(ctx context.Context, filter domain.ExpenseFilter) ([]domain.Expense, error)
	GetExpense(ctx context.Context, id int64) (domain.Expense, error)
	DeleteExpense(ctx context.Context, id int64) error
	PatchExpense(ctx context.Context, id int64, toPatch domain.Expense) (domain.Expense, error)
}

type AccountAccessChecker interface {
	HasAccess(ctx context.Context, uid, accountID int64) (bool, error)
	IsOwner(ctx context.Context, uid, accountID int64) (bool, error)
}

type CategoryAccessChecker interface {
	HasAccess(ctx context.Context, uid, categoryID, accountID int64) (bool, error)
}
