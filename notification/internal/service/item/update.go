package item

import (
	"context"

	"github.com/dfg007star/avito_informer/notification/internal/model"
)

func (s *service) UpdateItem(ctx context.Context, item *model.Item) error {
	err := s.itemRepository.UpdateItem(ctx, item)
	if err != nil {
		return err
	}

	return nil
}
