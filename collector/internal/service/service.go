package service

import (
	"context"

	"github.com/dfg007star/avito_informer/converter/internal/model"
)

type LinkService interface {
	GetAllLinks(ctx context.Context) ([]*model.Link, error)
}
