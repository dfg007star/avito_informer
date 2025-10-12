package auth

import (
	"net/http"

	"github.com/dfg007star/avito_informer/http/internal/config"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login" {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("auth")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		if cookie.Value != config.AppConfig().HTTP.Password() {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
