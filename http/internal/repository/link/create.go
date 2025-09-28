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

func (r *repository) CreateLink(ctx context.Context, link *model.Link) (*model.Link, error) {
	query, args, err := squirrel.Insert("links").
		Columns("name", "url").
		Values(link.Name, link.Url).
		Suffix("RETURNING *").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	row := r.db.QueryRow(ctx, query, args...)

	var createdLink repoModel.Link
	var updatedAt sql.NullTime
	err = row.Scan(
		&createdLink.ID,
		&createdLink.Name,
		&createdLink.Url,
		&createdLink.ParsedAt,
		&createdLink.CreatedAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	return converter.RepoLinkToModel(&createdLink), nil
}
