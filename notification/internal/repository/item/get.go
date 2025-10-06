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
		"id",
		"link_id",
		"uid",
		"title",
		"description",
		"url",
		"preview_url",
		"price",
		"need_notify",
		"created_at",
	).From("items").PlaceholderFormat(squirrel.Dollar).ToSql()
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
			&item.NeedNotify,
			&item.CreatedAt,
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
		"id",
		"link_id",
		"uid",
		"title",
		"description",
		"url",
		"preview_url",
		"price",
		"need_notify",
		"created_at",
	).From("items").Where(squirrel.Eq{"need_notify": false}).PlaceholderFormat(squirrel.Dollar).ToSql()
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
			&item.NeedNotify,
			&item.CreatedAt,
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
