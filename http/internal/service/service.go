package service

import (
	"context"
	"github.com/dfg007star/avito_informer/http/internal/model"
)

type LinkService interface {
	GetAllLinks(ctx context.Context) ([]*model.Link, error)
	GetLinkById(ctx context.Context, id string) (*model.Link, error)
	GetLinkItems(ctx context.Context, link *model.Link) ([]*model.Item, error)
}
