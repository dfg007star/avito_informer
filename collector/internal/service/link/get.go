package link

import (
	"context"

	"github.com/dfg007star/avito_informer/converter/internal/model"
)

func (s *service) GetAllLinks(ctx context.Context) ([]*model.Link, error) {
	links, err := s.linkRepository.GetAllLinks(ctx)
	if err != nil {
		return nil, err
	}

	return links, nil
}
