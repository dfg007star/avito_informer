package link

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/dfg007star/avito_informer/http/internal/model"
	"github.com/dfg007star/avito_informer/http/internal/repository/converter"
	repoModel "github.com/dfg007star/avito_informer/http/internal/repository/model"
)

func (r *repository) GetAllLinks(ctx context.Context) ([]*model.Link, error) {
	query, args, err := squirrel.Select("*").
		From("links").
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
	var updatedAt sql.NullTime
	for rows.Next() {
		var l repoModel.Link
		err := rows.Scan(
			&l.ID,
			&l.Name,
			&l.Url,
			&l.ParsedAt,
			&l.CreatedAt,
			&updatedAt,
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
	query, args, err := squirrel.Select("*").
		From("links").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	row := r.db.QueryRow(ctx, query, args...)

	var link repoModel.Link
	var updatedAt sql.NullTime
	err = row.Scan(
		&link.ID,
		&link.Name,
		&link.Url,
		&link.ParsedAt,
		&link.CreatedAt,
		&updatedAt,
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
		"need_notify",
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
			&item.NeedNotify,
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
