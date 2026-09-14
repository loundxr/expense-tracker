package accounts_transport_http

import (
	"time"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

type AccountDTOResponse struct {
	ID        int64     `json:"id"`
	Version   int64     `json:"version"`
	Name      string    `json:"name"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func accountDTOFromDomain(a domain.Account) AccountDTOResponse {
	return AccountDTOResponse{
		ID:        a.ID,
		Version:   a.Version,
		Name:      a.Name,
		UserID:    a.UserID,
		CreatedAt: a.CreatedAt,
	}
}

func accountDTOsFromDomains(accs []domain.Account) []AccountDTOResponse {
	res := make([]AccountDTOResponse, len(accs))
	for i, a := range accs {
		res[i] = accountDTOFromDomain(a)
	}
	return res
}
