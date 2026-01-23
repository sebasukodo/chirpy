package handler

import (
	"net/http"
)

func (cfg *ApiConfig) MiddlewareCheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		_, err := cfg.ValidateSessionID(w, r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) MiddlewareCheckAuthLoginPage(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		_, err := cfg.ValidateSessionID(w, r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return

	})
}
