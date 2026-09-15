package domain

import "time"

type SummaryFilter struct {
	AccountID int64
	From      *time.Time
	To        *time.Time
}

type MaxExpenseSummary struct {
	ID          int64
	Amount      int64
	Description string
	Date        time.Time
}

type TopCategorySummary struct {
	CategoryID   int64
	CategoryName string
	TotalAmount  int64
}

type Summary struct {
	AccountID         int64
	TotalAmount       int64
	TotalTransactions int64
	AverageAmount     int64
	MaxExpense        *MaxExpenseSummary
	TopCategory       *TopCategorySummary
}

func NewSummary(
	accID, totalAmount, totalTransactions, averageAmount int64,
	maxExpense *MaxExpenseSummary, topCategory *TopCategorySummary,
) Summary {
	return Summary{
		AccountID:         accID,
		TotalAmount:       totalAmount,
		TotalTransactions: totalTransactions,
		AverageAmount:     averageAmount,
		MaxExpense:        maxExpense,
		TopCategory:       topCategory,
	}
}

func NewUninitializedSummary(accountID int64) Summary {
	return NewSummary(accountID, uninitStatistics, uninitStatistics, uninitStatistics, nil, nil)
}

func NewMaxExpenseSummary(id int64, amount int64, description string, date time.Time) MaxExpenseSummary {
	return MaxExpenseSummary{
		ID:          id,
		Amount:      amount,
		Description: description,
		Date:        date,
	}
}

func NewSummaryFilter(accID int64, from, to *time.Time) SummaryFilter {
	return SummaryFilter{
		AccountID: accID,
		From:      from,
		To:        to,
	}
}
