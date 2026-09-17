package transport

import (
	"github.com/loundxr/expense-tracker/internal/core/domain"
	"github.com/loundxr/expense-tracker/internal/utils/money"
)

type MemberContributionDTO struct {
	UserID        int64   `json:"user_id"`
	Email         string  `json:"email"`
	TotalAmount   float64 `json:"total_amount"`
	ExpensesCount int64   `json:"expenses_count"`
	Percentage    float64 `json:"percentage"`
}

type MembersContributionDTOResponse struct {
	AccountID int64                   `json:"account_id"`
	Members   []MemberContributionDTO `json:"members"`
}

func membersContributionDTOFromDomains(accountID int64, members []domain.MemberContribution) MembersContributionDTOResponse {
	res := make([]MemberContributionDTO, len(members))
	for i, m := range members {
		res[i] = memberContributionDTOFromDomain(m)
	}
	return MembersContributionDTOResponse{
		AccountID: accountID,
		Members:   res,
	}
}

func memberContributionDTOFromDomain(m domain.MemberContribution) MemberContributionDTO {
	return MemberContributionDTO{
		UserID:        m.UserID,
		Email:         m.Email,
		TotalAmount:   money.ToUnits(m.TotalAmount),
		ExpensesCount: m.ExpensesCount,
		Percentage:    m.Percentage,
	}
}
