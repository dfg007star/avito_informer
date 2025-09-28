package link

import (
	"context"
	"github.com/dfg007star/avito_informer/http/internal/model"
)

func (s *service) CreateLink(ctx context.Context, link *model.Link) (*model.Link, error) {
	items, err := s.linkRepository.CreateLink(ctx, link)
	if err != nil {
		return nil, err
	}

	return items, nil
}
