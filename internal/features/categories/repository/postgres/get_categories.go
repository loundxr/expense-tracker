package categories_repository_postgres

import (
	"context"
	"fmt"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

func (r *CategoriesRepository) GetAllCategories(ctx context.Context) ([]domain.Category, error) {
	const op = "categories.repository.postgres.GetAllCategories"

	query := `
	SELECT id, version, name, user_id, created_at
	FROM expense_tracker.categories
	ORDER BY id ASC;`

	categories, err := r.fetch(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return categories, nil
}

func (r *CategoriesRepository) GetCategoriesByID(ctx context.Context, uid int) ([]domain.Category, error) {
	const op = "categories.repository.postgres.GetCategoriesByID"

	query := `
	SELECT id, version, name, user_id, created_at
	FROM expense_tracker.categories
	WHERE user_id IS NULL OR user_id=$1
	ORDER BY name ASC;`

	categories, err := r.fetch(ctx, query, uid)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return categories, nil
}

func (r *CategoriesRepository) fetch(ctx context.Context, sql string, args ...any) ([]domain.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []CategoryModel

	for rows.Next() {
		var cm CategoryModel
		if err := rows.Scan(
			&cm.ID,
			&cm.Version,
			&cm.Name,
			&cm.UserID,
			&cm.CreatedAt,
		); err != nil {
			return nil, err
		}

		models = append(models, cm)
	}

	res := categoryDomainsFromModels(models)
	return res, nil
}
