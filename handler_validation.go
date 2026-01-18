package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type parameters struct {
	Body string `json:"body"`
}

type returnError struct {
	Error string `json:"error"`
}

type returnValid struct {
	Valid bool `json:"valid"`
}

func handlerChirpValidation(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	respBody := returnValid{
		Valid: true,
	}

	respondWithJSON(w, 200, respBody)

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
