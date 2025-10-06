package repository

import (
	"context"

	"github.com/dfg007star/avito_informer/notification/internal/model"
)

type ItemRepository interface {
	GetItems(ctx context.Context) ([]*model.Item, error)
	GetNotNotifiedItems(ctx context.Context) ([]*model.Item, error)
	UpdateItem(ctx context.Context, item *model.Item) error
}
