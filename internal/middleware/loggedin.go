package middleware

import (
	"net/http"

	"github.com/sebasukodo/chirpy/internal/handler"
)

func CheckAuth(cfg *handler.ApiConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			next.ServeHTTP(w, r)
		})
	}
}
