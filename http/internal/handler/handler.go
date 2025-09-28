package handler

import (
	"net/http"
)

type LinkHandler interface {
	IndexHandler(w http.ResponseWriter, r *http.Request)
	CreateLinkHandler(w http.ResponseWriter, r *http.Request)
}
