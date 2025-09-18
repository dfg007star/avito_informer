package link

import (
	"html/template"

	"github.com/dfg007star/avito_informer/http/internal/service"
)

type handler struct {
	linkService service.LinkService
	templates   *template.Template
}

func NewApi(linkService service.LinkService, templates *template.Template) *handler {
	return &handler{
		linkService: linkService,
		templates:   templates,
	}
}
