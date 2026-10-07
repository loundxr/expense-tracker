package transport

import (
	"github.com/loundxr/expense-tracker/internal/core/domain"
	"github.com/loundxr/expense-tracker/internal/utils/money"
)

type CategoriesBreakdownDTOResponse struct {
	AccountID  int64                  `json:"account_id"`
	Categories []CategoryBreakdownDTO `json:"categories_breakdown"`
}

type CategoryBreakdownDTO struct {
	CategoryID        int64   `json:"category_id"`
	CategoryName      string  `json:"category_name"`
	TotalAmount       float64 `json:"total_amount"`
	TotalTransactions int64   `json:"total_transactions"`
	Percentage        float64 `json:"percentage"`
}

func categoriesBreakdownDTOFromDomains(b []domain.CategoryBreakdown) []CategoryBreakdownDTO {
	res := make([]CategoryBreakdownDTO, len(b))
	for i, c := range b {
		res[i] = categoryBreakdownDTOFromDomain(c)
	}
	return res
}

func categoryBreakdownDTOFromDomain(b domain.CategoryBreakdown) CategoryBreakdownDTO {
	return CategoryBreakdownDTO{
		CategoryID:        b.CategoryID,
		CategoryName:      b.CategoryName,
		TotalAmount:       money.ToUnits(b.TotalAmount),
		TotalTransactions: b.TotalTransactions,
		Percentage:        b.Percentage,
	}
}
