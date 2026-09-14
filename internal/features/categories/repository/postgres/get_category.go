package categories_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_postgres_pool "github.com/loundxr/expense-tracker/internal/core/repository/postgres/pool"
)

func (r *CategoriesRepository) GetCategory(ctx context.Context, id int64) (domain.Category, error) {
	const op = "categories.repository.postgres.GetCategory"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	SELECT id, version, name, user_id, created_at
	FROM expense_tracker.categories
	WHERE id=$1`

	row := r.pool.QueryRow(ctx, query, id)

	var cm CategoryModel

	err := row.Scan(
		&cm.ID,
		&cm.Version,
		&cm.Name,
		&cm.UserID,
		&cm.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNotFound) {
			return domain.Category{}, fmt.Errorf("%s: %w", op, core_errors.ErrNotFound)
		}
		return domain.Category{}, fmt.Errorf("%s: %w", op, err)
	}

	categoryDomain := categoryDomainFromModel(cm)
	return categoryDomain, nil
}

// TODO: fix category scope for account, not only for creator
func (r *CategoriesRepository) HasAccess(ctx context.Context, uid, categoryID int64) (bool, error) {
	const op = "categories.repository.postgres.HasAccess"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	SELECT EXISTS (
		SELECT 1 FROM expense_tracker.categories
		WHERE id=$1 AND (user_id=$2 OR user_id IS NULL)
	);`

	var exists bool
	err := r.pool.QueryRow(ctx, query, categoryID, uid).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return exists, nil
}
