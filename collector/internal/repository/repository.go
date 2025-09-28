package repository

import (
	"context"

	"github.com/dfg007star/avito_informer/collector/internal/model"
)

type LinkRepository interface {
	GetAllLinks(ctx context.Context) ([]*model.Link, error)
}
