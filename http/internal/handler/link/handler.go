package link

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

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
	fmt.Println(err)
	if err != nil {
		links = []*model.Link{}
	}
	fmt.Println(links)

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
	minPriceStr := r.FormValue("min_price")
	maxPriceStr := r.FormValue("max_price")

	var minPrice sql.NullInt64
	if minPriceStr != "" {
		val, err := strconv.ParseInt(minPriceStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid min_price", http.StatusBadRequest)
			return
		}
		minPrice = sql.NullInt64{Int64: val, Valid: true}
	}

	var maxPrice sql.NullInt64
	if maxPriceStr != "" {
		val, err := strconv.ParseInt(maxPriceStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid max_price", http.StatusBadRequest)
			return
		}
		maxPrice = sql.NullInt64{Int64: val, Valid: true}
	}

	link := &model.Link{
		Name:     name,
		Url:      linkUrl,
		MinPrice: minPrice,
		MaxPrice: maxPrice,
	}
	result, err := h.linkService.CreateLink(r.Context(), link)
	fmt.Println("handler create link", result, err)
	if err != nil {
		log.Fatalf("error creating link: %", err)
	}
	err = h.templates.ExecuteTemplate(w, "links_link", result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *handler) ShowLinkHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	link, err := h.linkService.GetLinkById(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if link == nil {
		http.NotFound(w, r)
	}
	//TODO: options its params from url
	//options := r.URL.Query()
	items, err := h.linkService.GetLinkItems(r.Context(), link)
	data := map[string]any{
		"ItemsLength": len(items),
		"LinkID":      id,
		"Items":       items,
	}
	err = h.templates.ExecuteTemplate(w, "link", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//
// func (h *Handler) ParseLinkHandler(w http.ResponseWriter, r *http.Request) {
// 		id := r.PathValue("id")
// 		link := h.service.GetLinkById(r.Context(), id)
// 		link.ParsedAt = sql.NullTime{time.Now().UTC(), true}
// 		updatedLink, err := h.service.UpdateLink(r.Context(), link)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 		go parser.ParseUrl(h.service.CreateEntitiesFromParsedData, updatedLink)
// 		err = h.templates.ExecuteTemplate(w, "links_link", extendLink(updatedLink))
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// }

func (h *handler) DeleteLinkHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.linkService.DeleteLink(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(200)
}
