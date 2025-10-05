package service

import (
	"context"

	"github.com/dfg007star/avito_informer/collector/internal/model"
	"github.com/dfg007star/avito_informer/collector/internal/repository"
	itemService "github.com/dfg007star/avito_informer/collector/internal/service/item"
	linkService "github.com/dfg007star/avito_informer/collector/internal/service/link"
)

type LinkService interface {
	GetAllLinks(ctx context.Context) ([]*model.Link, error)
}

type ItemService interface {
	CreateItems(ctx context.Context, items []*model.Item) error
}

type Service struct {
	LinkService
	ItemService
}

func New(repo *repository.Repository) *Service {
	return &Service{
		LinkService: linkService.NewLinkService(repo.LinkRepository),
		ItemService: itemService.NewItemService(repo.ItemRepository),
	}
}
