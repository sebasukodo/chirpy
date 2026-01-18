package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sebasukodo/chirpy/internal/database"
)

var slurs = [3]string{"kerfuffle", "sharbert", "fornax"}

type chirpCreateRequest struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type chirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirpCreate(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	chirpReq := chirpCreateRequest{}

	if err := decoder.Decode(&chirpReq); err != nil {
		respondWithError(w, 500, fmt.Sprintf("could not decode json message: %v", err))
		return
	}

	if len(chirpReq.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	chirpParam := database.CreateChirpParams{
		Body:   removeSlurs(chirpReq.Body),
		UserID: chirpReq.UserID,
	}

	data, err := cfg.dbQueries.CreateChirp(r.Context(), chirpParam)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("could not create chirp: %v", err))
		return
	}

	chirpResp := chirpResponse{
		ID:        data.ID,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		Body:      data.Body,
		UserID:    data.UserID,
	}

	respondWithJSON(w, 201, chirpResp)

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
