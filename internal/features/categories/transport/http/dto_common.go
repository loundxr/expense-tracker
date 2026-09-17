package transport

import (
	"time"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

type CategoryDTOResponse struct {
	ID        int64     `json:"id"`
	Version   int64     `json:"version"`
	Name      string    `json:"name"`
	UserID    *int64    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func categoryDTOFromDomain(c domain.Category) CategoryDTOResponse {
	return CategoryDTOResponse{
		ID:        c.ID,
		Version:   c.Version,
		Name:      c.Name,
		UserID:    c.UserID,
		CreatedAt: c.CreatedAt,
	}
}

func categoryDTOsFromDomains(cats []domain.Category) []CategoryDTOResponse {
	res := make([]CategoryDTOResponse, len(cats))
	for i, c := range cats {
		res[i] = categoryDTOFromDomain(c)
	}
	return res
}
