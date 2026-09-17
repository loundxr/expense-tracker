package transport

import (
	"time"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	"github.com/loundxr/expense-tracker/internal/utils/money"
)

type SummaryDTOResponse struct {
	AccountID         int64           `json:"account_id"`
	TotalAmount       float64         `json:"total_amount"`
	TotalTransactions int64           `json:"total_transactions"`
	AverageAmount     float64         `json:"average_amount"`
	MaxExpense        *MaxExpenseDTO  `json:"max_expense,omitempty"`
	TopCategory       *TopCategoryDTO `json:"top_category,omitempty"`
}

type MaxExpenseDTO struct {
	ID          int64     `json:"id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

type TopCategoryDTO struct {
	CategoryID   int64   `json:"category_id"`
	CategoryName string  `json:"category_name"`
	TotalAmount  float64 `json:"total_amount"`
}

func summaryDTOFromDomain(s domain.Summary) SummaryDTOResponse {
	dto := SummaryDTOResponse{
		AccountID:         s.AccountID,
		TotalAmount:       money.ToUnits(s.TotalAmount),
		TotalTransactions: s.TotalTransactions,
		AverageAmount:     money.ToUnits(s.AverageAmount),
	}

	if s.MaxExpense != nil {
		dto.MaxExpense = maxExpenseDTOFromDomain(s.MaxExpense)
	}
	if s.TopCategory != nil {
		dto.TopCategory = topCategoryDTOFromDomain(s.TopCategory)
	}
	return dto
}

func maxExpenseDTOFromDomain(me *domain.MaxExpenseSummary) *MaxExpenseDTO {
	return &MaxExpenseDTO{
		ID:          me.ID,
		Amount:      money.ToUnits(me.Amount),
		Description: me.Description,
		Date:        me.Date,
	}
}

func topCategoryDTOFromDomain(tc *domain.TopCategorySummary) *TopCategoryDTO {
	return &TopCategoryDTO{
		CategoryID:   tc.CategoryID,
		CategoryName: tc.CategoryName,
		TotalAmount:  money.ToUnits(tc.TotalAmount),
	}
}
