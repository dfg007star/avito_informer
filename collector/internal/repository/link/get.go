package link

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/dfg007star/avito_informer/collector/internal/model"
	"github.com/dfg007star/avito_informer/collector/internal/repository/converter"
	repoModel "github.com/dfg007star/avito_informer/collector/internal/repository/model"
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
