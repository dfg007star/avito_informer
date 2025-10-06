package link

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/dfg007star/avito_informer/http/internal/model"
	"github.com/dfg007star/avito_informer/http/internal/repository/converter"
	repoModel "github.com/dfg007star/avito_informer/http/internal/repository/model"
)

func (r *repository) GetAllLinks(ctx context.Context) ([]*model.Link, error) {
	query, args, err := squirrel.Select(
		"id",
		"name",
		"url",
		"min_price",
		"max_price",
		"parsed_at",
		"created_at",
	).From("links").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var links []*repoModel.Link
	for rows.Next() {
		var l repoModel.Link
		err := rows.Scan(
			&l.ID,
			&l.Name,
			&l.Url,
			&l.MinPrice,
			&l.MaxPrice,
			&l.ParsedAt,
			&l.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		links = append(links, &l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	var result []*model.Link
	for _, link := range links {
		result = append(result, converter.RepoLinkToModel(link))
	}

	return result, nil
}

func (r *repository) GetLinkById(ctx context.Context, id string) (*model.Link, error) {
	query, args, err := squirrel.Select(
		"id",
		"name",
		"url",
		"min_price",
		"max_price",
		"parsed_at",
		"created_at",
	).From("links").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	row := r.db.QueryRow(ctx, query, args...)

	var link repoModel.Link
	err = row.Scan(
		&link.ID,
		&link.Name,
		&link.Url,
		&link.MinPrice,
		&link.MaxPrice,
		&link.ParsedAt,
		&link.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	return converter.RepoLinkToModel(&link), nil
}

func (r *repository) GetLinkItems(ctx context.Context, link *model.Link) ([]*model.Item, error) {
	queryBuilder := squirrel.Select(
		"id",
		"link_id",
		"uid",
		"title",
		"description",
		"url",
		"preview_url",
		"price",
		"is_notify",
		"created_at",
	).
		From("items").
		Where(squirrel.Eq{"link_id": link.ID}).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var items []repoModel.Item
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
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	var result []*model.Item
	for _, item := range items {
		result = append(result, converter.RepoItemToModel(&item))
	}

	return result, nil
}
