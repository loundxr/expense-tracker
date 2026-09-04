package accounts_repository_postgres

import (
	"time"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

type AccountModel struct {
	ID        int
	Version   int
	Name      string
	UserID    int
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
