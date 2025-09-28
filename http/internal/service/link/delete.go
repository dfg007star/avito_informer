package link

import (
	"context"
)

func (s *service) DeleteLink(ctx context.Context, id string) error {
	err := s.linkRepository.DeleteLink(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
