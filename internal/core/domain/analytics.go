package domain

import "time"

// Summary entity with constructors
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

// SummaryFilter to create Get Summary request
type SummaryFilter struct {
	AccountID int64
	From      *time.Time
	To        *time.Time
}

func NewSummaryFilter(accID int64, from, to *time.Time) SummaryFilter {
	return SummaryFilter{
		AccountID: accID,
		From:      from,
		To:        to,
	}
}

// Summary fields

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

// CategoriesBreakdown entity
type CategoryBreakdown struct {
	CategoryID        int64
	CategoryName      string
	TotalAmount       int64
	TotalTransactions int64
	Percentage        float64
}

// ExpenseTrendPoint entity
type ExpenseTrendPoint struct {
	Date          time.Time
	TotalAmount   int64
	ExpensesCount int64
}

// Trends filter to create Get Trends request
type TrendsFilter struct {
	AccountID int64
	From      *time.Time
	To        *time.Time
	Interval  string
}

func NewTrendsFilter(accID int64, from, to *time.Time, interval string) TrendsFilter {
	return TrendsFilter{
		AccountID: accID,
		From:      from,
		To:        to,
		Interval:  interval,
	}
}

// MemberContribution entity with
type MemberContribution struct {
	UserID        int64
	Email         string
	TotalAmount   int64
	ExpensesCount int64
	Percentage    float64
}
