package repository

import (
	"time"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

type AccountModel struct {
	ID        int64
	Version   int64
	Name      string
	UserID    int64
	CreatedAt time.Time
}

func accountDomainFromModel(am AccountModel) domain.Account {
	return domain.NewAccount(
		am.ID,
		am.Version,
		am.Name,
		am.UserID,
		am.CreatedAt,
	)
}

func accountDomainsFromModels(ams []AccountModel) []domain.Account {
	domains := make([]domain.Account, len(ams))
	for i, m := range ams {
		domains[i] = accountDomainFromModel(m)
	}
	return domains
}
