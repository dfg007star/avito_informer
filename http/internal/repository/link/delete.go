package link

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
)

func (r *repository) DeleteLink(ctx context.Context, id string) error {

	query, args, err := squirrel.Delete("links").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}
