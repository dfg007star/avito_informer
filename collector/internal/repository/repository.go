package repository

import (
	"context"

	"github.com/dfg007star/avito_informer/collector/internal/model"
	"github.com/dfg007star/avito_informer/collector/internal/repository/item"
	"github.com/dfg007star/avito_informer/collector/internal/repository/link"
	"github.com/jackc/pgx/v5"
)

type LinkRepository interface {
	GetAllLinks(ctx context.Context) ([]*model.Link, error)
}

type ItemRepository interface {
	CreateItems(ctx context.Context, items []*model.Item) error
}

type Repository struct {
	LinkRepository
	ItemRepository
}

func New(clientPostgres *pgx.Conn) *Repository {
	return &Repository{
		LinkRepository: link.NewRepository(clientPostgres),
		ItemRepository: item.NewRepository(clientPostgres),
	}
}
