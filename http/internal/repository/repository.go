package repository

import (
	"context"
	"github.com/dfg007star/avito_informer/http/internal/model"
)

type LinkRepository interface {
	GetAllLinks(ctx context.Context, orderUuid string) ([]*model.Link, error)
}
