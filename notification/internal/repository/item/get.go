package item

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/dfg007star/avito_informer/notification/internal/model"
	"github.com/dfg007star/avito_informer/notification/internal/repository/converter"
	repoModel "github.com/dfg007star/avito_informer/notification/internal/repository/model"
)

func (r *repository) GetItems(ctx context.Context) ([]*model.Item, error) {
	query, args, err := squirrel.Select(
		"items.id",
		"items.link_id",
		"items.uid",
		"items.title",
		"items.description",
		"items.url",
		"items.preview_url",
		"items.price",
		"items.is_notify",
		"items.created_at",
		"links.name AS category_title",
		"links.url AS link_url",
	).From("items").
		LeftJoin("links ON items.link_id = links.id").
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var items []*repoModel.Item
	for rows.Next() {
		var item repoModel.Item
		err := rows.Scan(
			&item.ID,
			&item.LinkId,
			&item.Uid,
			&item.Title,
			&item.Description,
			&item.Url,
			&item.PreviewUrl,
			&item.Price,
			&item.IsNotify,
			&item.CreatedAt,
			&item.CategoryTitle,
			&item.LinkUrl,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return converter.RepoItemToModelList(items), nil
}

func (r *repository) GetNotNotifiedItems(ctx context.Context) ([]*model.Item, error) {
	query, args, err := squirrel.Select(
		"items.id",
		"items.link_id",
		"items.uid",
		"items.title",
		"items.description",
		"items.url",
		"items.preview_url",
		"items.price",
		"items.is_notify",
		"items.created_at",
		"links.name AS category_title",
		"links.url AS link_url",
	).From("items").
		LeftJoin("links ON items.link_id = links.id").
		Where(squirrel.Eq{"items.is_notify": false}).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var items []*repoModel.Item
	for rows.Next() {
		var item repoModel.Item
		err := rows.Scan(
			&item.ID,
			&item.LinkId,
			&item.Uid,
			&item.Title,
			&item.Description,
			&item.Url,
			&item.PreviewUrl,
			&item.Price,
			&item.IsNotify,
			&item.CreatedAt,
			&item.CategoryTitle,
			&item.LinkUrl,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return converter.RepoItemToModelList(items), nil
}
