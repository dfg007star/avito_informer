package item

import (
	"context"

	"github.com/dfg007star/avito_informer/notification/internal/model"
)

func (s *service) GetItems(ctx context.Context) ([]*model.Item, error) {
	items, err := s.itemRepository.GetItems(ctx)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *service) GetNotNotifiedItems(ctx context.Context) ([]*model.Item, error) {
	items, err := s.itemRepository.GetNotNotifiedItems(ctx)
	if err != nil {
		return nil, err
	}

	return items, nil
}
