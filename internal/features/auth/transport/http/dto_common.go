package transport

import (
	"time"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

type UserDTOResponse struct {
	ID        int64     `json:"id"`
	Version   int64     `json:"version"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	Role      string    `json:"role"`
}

func userDTOFromDomain(u domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:        u.ID,
		Version:   u.Version,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		Role:      u.Role,
	}
}
