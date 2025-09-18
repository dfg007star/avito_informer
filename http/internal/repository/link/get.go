package link

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/dfg007star/avito_informer/internal/model"
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

	var links []*model.Link
	for rows.Next() {
		var l model.Link
		// Replace these fields with actual fields from your model.Link
		err := rows.Scan(
			&l.Name,
			&l.Url,
			&l.CreatedAt,
			&l.ParsedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		links = append(links, &l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return links, nil
}
