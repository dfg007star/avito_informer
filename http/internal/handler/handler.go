package handler

import (
	"net/http"
)

type LinkHandler interface {
	IndexHandler(w http.ResponseWriter, r *http.Request)
	CreateLinkHandler(w http.ResponseWriter, r *http.Request)
	ShowLinkHandler(w http.ResponseWriter, r *http.Request)
	DeleteLinkHandler(w http.ResponseWriter, r *http.Request)
}

type AuthHandler interface {
	LoginHandler(w http.ResponseWriter, r *http.Request)
}

type LogHandler interface {
	IndexHandler(w http.ResponseWriter, r *http.Request)
}
