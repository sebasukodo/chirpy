package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {

	if cfg.platform != "dev" {
		respondWithError(w, 403, "ACCESS DENIED")
		return
	}

	if err := cfg.dbQueries.DeleteAllUsers(r.Context()); err != nil {
		respondWithError(w, 500, fmt.Sprintf("could not delete all users: %w", err))
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	body := "<html><body><h1>All Users have been deleted successfully</h1></body></html>"
	w.Write([]byte(body))

}
