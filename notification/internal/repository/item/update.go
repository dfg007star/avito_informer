package item

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/dfg007star/avito_informer/notification/internal/model"
)

func (r *repository) UpdateItem(ctx context.Context, item *model.Item) error {
	query, args, err := squirrel.Update("items").
		Set("need_notify", item.NeedNotify).
		Where(squirrel.Eq{"id": item.ID}).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}

	return nil
}
