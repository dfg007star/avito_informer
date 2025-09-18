package link

import (
	"context"

	"github.com/dfg007star/avito_informer/http/internal/model"
)

func (s *service) GetAllLinks(ctx context.Context) ([]*model.Link, error) {
	links, err := s.linkRepository.GetAllLinks(ctx)
	if err != nil {
		return nil, err
	}

	return links, nil
}

func (s *service) GetLinkById(ctx context.Context, id string) (*model.Link, error) {
	link, err := s.linkRepository.GetLinkById(ctx, id)
	if err != nil {
		return nil, err
	}

	return link, nil
}

func (s *service) GetLinkItems(ctx context.Context, link *model.Link) ([]*model.Item, error) {
	items, err := s.linkRepository.GetLinkItems(ctx, link)
	if err != nil {
		return nil, err
	}

	return items, nil
}
