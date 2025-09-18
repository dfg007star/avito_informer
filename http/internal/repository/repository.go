package repository

import (
	"context"

	"github.com/dfg007star/avito_informer/internal/model"
)

type OrderRepository interface {
	GetAllLinks(ctx context.Context, orderUuid string) ([]*model.Link, error)
}
