package repository

import (
	"time"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

type UserModel struct {
	ID           int64
	Version      int64
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	Role         string
}

func userDomainFromModel(um UserModel) domain.User {
	return domain.NewUser(
		um.ID,
		um.Version,
		um.Email,
		um.PasswordHash,
		um.CreatedAt,
		um.Role,
	)
}

func userDomainsFromModels(ums []UserModel) []domain.User {
	domains := make([]domain.User, len(ums))
	for i, um := range ums {
		domains[i] = userDomainFromModel(um)
	}
	return domains
}
