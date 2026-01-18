package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var slurs = [3]string{"kerfuffle", "sharbert", "fornax"}

type parameters struct {
	Body string `json:"body"`
}

type returnError struct {
	Error string `json:"error"`
}

type returnValid struct {
	Valid bool `json:"valid"`
}

type returnBody struct {
	CleanedBody string `json:"cleaned_body"`
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

	respBody := returnBody{
		CleanedBody: removeSlurs(params.Body),
	}

	respondWithJSON(w, 200, respBody)

}

func removeSlurs(msg string) string {

	splittedMsg := strings.Split(msg, " ")

	for _, slur := range slurs {

		for i := range splittedMsg {

			if strings.ToLower(splittedMsg[i]) == slur {
				splittedMsg[i] = "****"
			}
		}

	}

	return strings.Join(splittedMsg, " ")

}
