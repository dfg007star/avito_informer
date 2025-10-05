package item

import (
	"github.com/dfg007star/avito_informer/collector/internal/repository"
)

type service struct {
	itemRepository repository.ItemRepository
}

func NewItemService(
	itemRepository repository.ItemRepository,
) *service {
	return &service{
		itemRepository: itemRepository,
	}
}
