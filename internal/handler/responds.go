package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type returnError struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {

	respBody := returnError{
		Error: msg,
	}

	respondWithJSON(w, code, respBody)

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	w.Header().Set("Content-Type", "application/json")

	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Write(dat)

}

func Readiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *ApiConfig) Reset(w http.ResponseWriter, r *http.Request) {

	if cfg.Platform != "dev" {
		respondWithError(w, 403, "ACCESS DENIED")
		return
	}

	if err := cfg.DbQueries.DeleteAllUsers(r.Context()); err != nil {
		respondWithError(w, 500, fmt.Sprintf("could not delete all users: %v", err))
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	body := "<html><body><h1>All Users have been deleted successfully</h1></body></html>"
	w.Write([]byte(body))

}
