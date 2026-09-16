package domain

import "time"

type Summary struct {
	AccountID         int64
	TotalAmount       int64
	TotalTransactions int64
	AverageAmount     int64
	MaxExpense        *MaxExpenseSummary
	TopCategory       *TopCategorySummary
}

type CategoryBreakdown struct {
	CategoryID        int64
	CategoryName      string
	TotalAmount       int64
	TotalTransactions int64
	Percentage        float64
}
type SummaryFilter struct {
	AccountID int64
	From      *time.Time
	To        *time.Time
}

type TrendsFilter struct {
	AccountID int64
	From      *time.Time
	To        *time.Time
	Interval  string
}

type ExpenseTrendPoint struct {
	Date          time.Time
	TotalAmount   int64
	ExpensesCount int64
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

func NewCategoriesBreakdown(
	categoryID int64,
	categoryName string,
	totalAmount int64,
	totalTransactions int64,
	percentage float64,
) CategoryBreakdown {
	return CategoryBreakdown{
		CategoryID:        categoryID,
		CategoryName:      categoryName,
		TotalAmount:       totalAmount,
		TotalTransactions: totalTransactions,
		Percentage:        percentage,
	}
}

func NewMaxExpenseSummary(id int64, amount int64, description string, date time.Time) MaxExpenseSummary {
	return MaxExpenseSummary{
		ID:          id,
		Amount:      amount,
		Description: description,
		Date:        date,
	}
}

func NewExpenseTrendPoint(date time.Time, totalAmount, ExpenseCount int64) ExpenseTrendPoint {
	return ExpenseTrendPoint{
		Date:          date,
		TotalAmount:   totalAmount,
		ExpensesCount: ExpenseCount,
	}
}

func NewSummaryFilter(accID int64, from, to *time.Time) SummaryFilter {
	return SummaryFilter{
		AccountID: accID,
		From:      from,
		To:        to,
	}
}

func NewTrendsFilter(accID int64, from, to *time.Time, interval string) TrendsFilter {
	return TrendsFilter{
		AccountID: accID,
		From:      from,
		To:        to,
		Interval:  interval,
	}
}
