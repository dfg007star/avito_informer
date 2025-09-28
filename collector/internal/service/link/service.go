package link

import (
	"github.com/dfg007star/avito_informer/converter/internal/repository"
)

type service struct {
	linkRepository repository.LinkRepository
}

func NewLinkService(
	linkRepository repository.LinkRepository,
) *service {
	return &service{
		linkRepository: linkRepository,
	}
}
