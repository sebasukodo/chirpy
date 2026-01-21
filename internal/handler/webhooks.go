package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/sebasukodo/chirpy/internal/auth"
)

type data struct {
	UserID uuid.UUID `json:"user_id"`
}

type eventRequest struct {
	Event string `json:"event"`
	Data  data   `json:"data"`
}

func (cfg *ApiConfig) VIP(w http.ResponseWriter, r *http.Request) {

	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "Access Denied")
		return
	}

	if key != cfg.PolkaApiKey {
		respondWithError(w, 401, "Access Denied")
		return
	}

	decoder := json.NewDecoder(r.Body)

	eventRequestData := eventRequest{}

	if err := decoder.Decode(&eventRequestData); err != nil {
		respondWithError(w, 500, "Internal Error")
		return
	}

	if eventRequestData.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	_, err = cfg.DbQueries.UpdateUserVIP(r.Context(), eventRequestData.Data.UserID)
	if err != nil {
		respondWithError(w, 404, "User not found")
		return
	}

	w.WriteHeader(204)

}
