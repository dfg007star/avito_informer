package service

import (
	"context"

	"github.com/dfg007star/avito_informer/collector/internal/model"
)

type LinkService interface {
	GetAllLinks(ctx context.Context) ([]*model.Link, error)
}
