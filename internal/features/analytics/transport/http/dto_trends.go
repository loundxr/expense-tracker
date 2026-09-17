package analytics_transport_http

import (
	"github.com/loundxr/expense-tracker/internal/core/domain"
	"github.com/loundxr/expense-tracker/internal/utils/money"
)

type TrendPointDTOResponse struct {
	Date          string  `json:"date"`
	TotalAmount   float64 `json:"total_amount"`
	ExpensesCount int64   `json:"expenses_count"`
}

type TrendsDTOResponse struct {
	AccountID int64                   `json:"account_id"`
	Interval  string                  `json:"interval"`
	Trends    []TrendPointDTOResponse `json:"trends"`
}

func trendsDTOResponseFromDomains(
	accountID int64,
	interval string,
	trends []domain.ExpenseTrendPoint,
) TrendsDTOResponse {
	dateFormat := "2006-01-02"
	if interval == "month" {
		dateFormat = "2006-01"
	}
	return TrendsDTOResponse{
		AccountID: accountID,
		Interval:  interval,
		Trends:    trendPointsDTOFromDomains(trends, dateFormat),
	}
}

func trendPointsDTOFromDomains(trends []domain.ExpenseTrendPoint, dateFormat string) []TrendPointDTOResponse {
	res := make([]TrendPointDTOResponse, len(trends))
	for i, t := range trends {
		res[i] = trendPointDTOFromDomain(t, dateFormat)
	}
	return res
}

func trendPointDTOFromDomain(t domain.ExpenseTrendPoint, dateFormat string) TrendPointDTOResponse {
	dateStr := t.Date.Format(dateFormat)
	return TrendPointDTOResponse{
		Date:          dateStr,
		TotalAmount:   money.ToUnits(t.TotalAmount),
		ExpensesCount: t.ExpensesCount,
	}
}
