package link

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/dfg007star/avito_informer/http/internal/model"
	"github.com/dfg007star/avito_informer/http/internal/repository/converter"
	repoModel "github.com/dfg007star/avito_informer/http/internal/repository/model"
)

func (r *repository) CreateLink(ctx context.Context, link *model.Link) (*model.Link, error) {
	query, args, err := squirrel.Insert("links").
		Columns("name", "url", "min_price", "max_price").
		Values(link.Name, link.Url, link.MinPrice, link.MaxPrice).
		Suffix("RETURNING id, name, url, min_price, max_price, parsed_at, created_at").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	row := r.db.QueryRow(ctx, query, args...)

	var createdLink repoModel.Link
	err = row.Scan(
		&createdLink.ID,
		&createdLink.Name,
		&createdLink.Url,
		&createdLink.MinPrice,
		&createdLink.MaxPrice,
		&createdLink.ParsedAt,
		&createdLink.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	return converter.RepoLinkToModel(&createdLink), nil
}
