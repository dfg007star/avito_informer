package service

import (
	"context"

	"github.com/dfg007star/avito_informer/notification/internal/model"
)

type ItemService interface {
	GetItems(ctx context.Context) ([]*model.Item, error)
	UpdateItem(ctx context.Context, item *model.Item) error
}

type TelegramService interface {
	SendItemNotification(ctx context.Context, event model.ItemEvent) error
}
