package categories_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_postgres_pool "github.com/loundxr/expense-tracker/internal/core/repository/postgres/pool"
)

func (r *CategoriesRepository) CreateCategory(
	ctx context.Context,
	c domain.Category,
) (domain.Category, error) {
	const op = "categories.repository.postgres.CreateCategory"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	INSERT INTO expense_tracker.categories (name, user_id, created_at)
	VALUES ($1, $2, $3)
	RETURNING id, version, name, user_id, created_at;`

	row := r.pool.QueryRow(ctx, query, c.Name, c.UserID, c.CreatedAt)

	var cm CategoryModel
	err := row.Scan(
		&cm.ID,
		&cm.Version,
		&cm.Name,
		&cm.UserID,
		&cm.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return domain.Category{}, fmt.Errorf(
				"%s: scan from returned row: %w",
				op,
				core_errors.ErrConflict,
			)
		}
		return domain.Category{}, fmt.Errorf("%s: scan from returned row: %w", op, err)
	}

	categoryDomain := categoryDomainFromModel(cm)
	return categoryDomain, nil
}
