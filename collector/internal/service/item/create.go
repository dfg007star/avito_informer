package item

import (
	"context"

	"github.com/dfg007star/avito_informer/collector/internal/model"
)

func (s *service) CreateItems(ctx context.Context, items []*model.Item) error {
	err := s.itemRepository.CreateItems(ctx, items)
	if err != nil {
		return err
	}

	return nil
}
