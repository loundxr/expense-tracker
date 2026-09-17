package repository

import (
	"time"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

type ExpenseModel struct {
	ID          int64
	Version     int64
	AccountID   int64
	UserID      int64
	CategoryID  int64
	Amount      int64
	Currency    string
	Description string
	Date        time.Time
	CreatedAt   time.Time
}

func expenseDomainFromModel(em ExpenseModel) domain.Expense {
	return domain.NewExpense(
		em.ID,
		em.Version,
		em.AccountID,
		em.UserID,
		em.CategoryID,
		em.Amount,
		em.Currency,
		em.Description,
		em.Date,
		em.CreatedAt,
	)
}

func expenseDomainsFromModels(models []ExpenseModel) []domain.Expense {
	res := make([]domain.Expense, len(models))
	for i, m := range models {
		res[i] = expenseDomainFromModel(m)
	}
	return res
}
