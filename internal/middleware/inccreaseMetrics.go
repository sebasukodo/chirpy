package middleware

import (
	"net/http"

	"github.com/sebasukodo/chirpy/internal/handler"
)

func MetricsInc(cfg *handler.ApiConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cfg.FileserverHits.Add(1)
			next.ServeHTTP(w, r)
		})
	}
}
