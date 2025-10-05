package item

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/dfg007star/avito_informer/collector/internal/model"
)

func (r *repository) CreateItems(ctx context.Context, items []*model.Item) error {
	if len(items) == 0 {
		return nil
	}

	fmt.Printf("creating %d items\n", len(items))

	queryBuilder := squirrel.Insert("items").
		Columns("link_id", "title", "price", "description", "url")

	for _, item := range items {
		queryBuilder = queryBuilder.Values(item.LinkId, item.Title, item.Price, item.Description, item.Url)
	}

	query, args, err := queryBuilder.
		Suffix("ON CONFLICT (url) DO NOTHING").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to exec query: %w", err)
	}

	return nil
}
