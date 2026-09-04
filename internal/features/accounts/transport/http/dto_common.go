package accounts_transport_http

import (
	"time"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

type AccountDTOResponse struct {
	ID        int       `json:"id"`
	Version   int       `json:"version"`
	Name      string    `json:"name"`
	UserID    int       `json:"user_id"`
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
