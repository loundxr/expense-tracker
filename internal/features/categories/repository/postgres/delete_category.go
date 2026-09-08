package categories_repository_postgres

import (
	"context"
	"fmt"
)

func (r *CategoriesRepository) DeleteCategory(ctx context.Context, id int64) error {
	const op = "categories.repository.postgres.DeleteCategory"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	DELETE FROM expense_tracker.categories
	WHERE id=$1;`

	if _, err := r.pool.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
