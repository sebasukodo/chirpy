package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sebasukodo/chirpy/internal/auth"
	"github.com/sebasukodo/chirpy/internal/database"
)

var slurs = [3]string{"kerfuffle", "sharbert", "fornax"}

type chirpCreateRequest struct {
	Body string `json:"body"`
}

type chirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *ApiConfig) ChirpsGetAll(w http.ResponseWriter, r *http.Request) {

	queryAuthor := r.URL.Query().Get("author_id")
	querySort := r.URL.Query().Get("sort")

	if querySort != "desc" {
		querySort = "asc"
	}

	var chirps []database.Chirp
	if queryAuthor == "" {
		var err error
		chirps, err = cfg.DbQueries.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, 500, fmt.Sprintf("could not retrieve all chirps: %v", err))
			return
		}

	} else {

		query, err := uuid.Parse(queryAuthor)
		if err != nil {
			respondWithError(w, 400, "invalid author_id")
			return
		}

		chirps, err = cfg.DbQueries.GetAllChirpsFromAuthor(r.Context(), query)
		if err != nil {
			respondWithError(w, 500, "could not retrieve chirps")
			return
		}

	}

	sort.Slice(chirps, func(i int, j int) bool {
		if querySort == "desc" {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		}
		return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
	})

	sortedChirps := make([]chirpResponse, 0, len(chirps))
	for _, chirp := range chirps {

		sortedChirps = append(sortedChirps, convertDatabaseChirp(chirp))

	}

	respondWithJSON(w, 200, sortedChirps)

}

func (cfg *ApiConfig) ChirpsGetByID(w http.ResponseWriter, r *http.Request) {

	chirpIDString := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("%v is not a valid uuid: %v", chirpIDString, err))
		return
	}

	chirp, err := cfg.DbQueries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, "chirp not found")
		return
	}

	respondWithJSON(w, 200, convertDatabaseChirp(chirp))

}

func (cfg *ApiConfig) ChirpsCreate(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	chirpReq := chirpCreateRequest{}

	if err := decoder.Decode(&chirpReq); err != nil {
		respondWithError(w, 500, fmt.Sprintf("could not decode json message: %v", err))
		return
	}

	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Access Denied")
		return
	}

	uid, err := auth.ValidateJWT(bearer, cfg.TokenSecret)
	if err != nil {
		respondWithError(w, 401, "Access Denied")
		return
	}

	if len(chirpReq.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	chirpParam := database.CreateChirpParams{
		Body:   removeSlurs(chirpReq.Body),
		UserID: uid,
	}

	data, err := cfg.DbQueries.CreateChirp(r.Context(), chirpParam)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("could not create chirp: %v", err))
		return
	}

	respondWithJSON(w, 201, convertDatabaseChirp(data))

}

func (cfg *ApiConfig) ChirpsDeleteByID(w http.ResponseWriter, r *http.Request) {

	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Access Denied")
		return
	}

	userID, err := auth.ValidateJWT(bearer, cfg.TokenSecret)
	if err != nil {
		respondWithError(w, 401, "Access Denied")
		return
	}

	chirpIDString := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, 400, "Access Denied")
		return
	}

	chirpUserID, err := cfg.DbQueries.GetChirpUserID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, "not found in database")
		return
	}

	if userID != chirpUserID {
		respondWithError(w, 403, "Access Denied")
		return
	}

	if err := cfg.DbQueries.DeleteChirpByID(r.Context(), chirpID); err != nil {
		respondWithError(w, 500, "could not delete chirp")
		return
	}

	w.WriteHeader(http.StatusNoContent)

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

func convertDatabaseChirp(dbChirp database.Chirp) chirpResponse {
	return chirpResponse{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
}
