package users_transport_http

import (
	"time"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

type UserDTOResponse struct {
	ID        int       `json:"id"`
	Version   int       `json:"version"`
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

func usersDTOFromDomains(users []domain.User) []UserDTOResponse {
	usersDTO := make([]UserDTOResponse, len(users))
	for i, u := range users {
		usersDTO[i] = userDTOFromDomain(u)
	}
	return usersDTO
}
