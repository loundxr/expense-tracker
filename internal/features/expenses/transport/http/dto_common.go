package transport

import (
	"time"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

type ExpenseDTOResponse struct {
	ID          int64     `json:"id"`
	Version     int64     `json:"version"`
	AccountID   int64     `json:"account_id"`
	UserID      int64     `json:"user_id"`
	CategoryID  int64     `json:"category_id"`
	Amount      int64     `json:"amount"`
	Currency    string    `json:"currency"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
}

func expenseDTOFromDomain(e domain.Expense) ExpenseDTOResponse {
	return ExpenseDTOResponse{
		ID:          e.ID,
		Version:     e.Version,
		AccountID:   e.AccountID,
		UserID:      e.UserID,
		CategoryID:  e.CategoryID,
		Amount:      e.Amount,
		Currency:    e.Currency,
		Description: e.Description,
		Date:        e.Date,
		CreatedAt:   e.CreatedAt,
	}
}

func expenseDTOsFromDomains(exps []domain.Expense) []ExpenseDTOResponse {
	res := make([]ExpenseDTOResponse, len(exps))
	for i, e := range exps {
		res[i] = expenseDTOFromDomain(e)
	}
	return res
}
