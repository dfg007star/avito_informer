package link

import (
	"html/template"
	"log"
	"net/http"

	"github.com/dfg007star/avito_informer/http/internal/model"
	"github.com/dfg007star/avito_informer/http/internal/service"
)

type handler struct {
	linkService service.LinkService
	templates   *template.Template
}

func NewHandler(linkService service.LinkService, templates *template.Template) *handler {
	return &handler{
		linkService: linkService,
		templates:   templates,
	}
}

func (h *handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	links, err := h.linkService.GetAllLinks(r.Context())
	if err != nil {
		links = []*model.Link{}
	}

	data := map[string]any{
		"Links": links,
	}
	err = h.templates.ExecuteTemplate(w, "index", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *handler) CreateLinkHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	linkUrl := r.FormValue("url")
	link := &model.Link{
		Name: name,
		Url:  linkUrl,
	}
	_, err := h.linkService.CreateLink(r.Context(), link)
	if err != nil {
		log.Fatal("<CreateHandler> Error creating link!")
	}
	err = h.templates.ExecuteTemplate(w, "links_link", extendLink(link))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// func (h *Handler) ShowLinkHandler(w http.ResponseWriter, r *http.Request) {
// 	id := r.PathValue("id")
// 	link := h.service.GetLinkById(r.Context(), id)
// 	if link == nil {
// 		http.NotFound(w, r)
// 	}
// 	options := r.URL.Query()
// 	items := h.service.GetLinkItems(r.Context(), link, &options)
// 	data := map[string]any{
// 		"ItemsLength": len(items),
// 		"LinkID":      id,
// 		"Items":       items,
// 	}
// 	err := h.templates.ExecuteTemplate(w, "link", data)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }
//
// func (h *Handler) ParseLinkHandler(w http.ResponseWriter, r *http.Request) {
// 	id := r.PathValue("id")
// 	link := h.service.GetLinkById(r.Context(), id)
// 	link.ParsedAt = sql.NullTime{time.Now().UTC(), true}
// 	updatedLink, err := h.service.UpdateLink(r.Context(), link)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// 	go parser.ParseUrl(h.service.CreateEntitiesFromParsedData, updatedLink)
// 	err = h.templates.ExecuteTemplate(w, "links_link", extendLink(updatedLink))
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }
//
// func (h *Handler) DeleteLinkHandler(w http.ResponseWriter, r *http.Request) {
// 	id := r.PathValue("id")
// 	err := h.service.DeleteLink(r.Context(), id)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// 	w.WriteHeader(200)
// }
