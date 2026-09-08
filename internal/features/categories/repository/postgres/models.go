package categories_repository_postgres

import (
	"time"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

type CategoryModel struct {
	ID        int
	Version   int
	Name      string
	UserID    *int
	CreatedAt time.Time
}

func categoryDomainFromModel(m CategoryModel) domain.Category {
	return domain.NewCategory(
		m.ID,
		m.Version,
		m.Name,
		m.UserID,
		m.CreatedAt,
	)
}

func categoryDomainsFromModels(models []CategoryModel) []domain.Category {
	res := make([]domain.Category, len(models))
	for i, m := range models {
		res[i] = categoryDomainFromModel(m)
	}
	return res
}
