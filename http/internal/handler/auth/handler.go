package auth

import (
	"html/template"
	"net/http"

	"github.com/dfg007star/avito_informer/http/internal/config"
)

type Handler struct {
	tmpl *template.Template
}

func NewHandler(tmpl *template.Template) *Handler {
	return &Handler{tmpl: tmpl}
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := h.tmpl.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if r.Method == http.MethodPost {
		password := r.FormValue("password")
		if password == config.AppConfig().HTTP.Password() {
			cookie := &http.Cookie{
				Name:  "auth",
				Value: password,
			}
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			err := h.tmpl.ExecuteTemplate(w, "login.html", "invalid password")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
