package handler

import (
	"net/http"
	"time"

	"github.com/sebasukodo/chirpy/internal/auth"
	"github.com/sebasukodo/chirpy/internal/database"
	"github.com/sebasukodo/chirpy/templates"
)

const SessionIDExpiresInHours = time.Duration(60*24) * time.Hour

func (cfg *ApiConfig) RefreshSessionID(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		respondWithError(w, r, 401, "Access Denied")
		return
	}

	sessionID, err := cfg.DbQueries.GetSessionIDByID(r.Context(), cookie.Value)
	if err != nil {
		respondWithError(w, r, 401, "Access Denied")
		return
	}

	if sessionID.RevokedAt.Valid {
		respondWithError(w, r, 401, "Access Denied")
		return
	}

	_, err = cfg.DbQueries.SetSessionIDInvalid(r.Context(), sessionID.ID)
	if err != nil {
		respondWithError(w, r, 500, "could not update database")
		return
	}

	newSessionID, err := auth.MakeSessionID()
	if err != nil {
		respondWithError(w, r, 500, "Could not refresh session")
		return
	}

	_, err = cfg.DbQueries.StoreSessionID(r.Context(), database.StoreSessionIDParams{
		ID:        newSessionID,
		ExpiresAt: time.Now().UTC().Add(SessionIDExpiresInHours),
	})
	if err != nil {
		respondWithError(w, r, 500, "Could not store new session")
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    newSessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().UTC().Add(SessionIDExpiresInHours),
	})

	respondWithHTML(templates.Layout(templates.HomepageContent(), "Success"), w, r)

}

func (cfg *ApiConfig) RevokeSessionID(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("session_id")
	if err != nil {
		respondWithError(w, r, 401, "Access Denied")
		return
	}

	sessionID, err := cfg.DbQueries.GetSessionIDByID(r.Context(), cookie.Value)
	if err != nil {
		respondWithError(w, r, 401, "Access Denied")
		return
	}

	_, err = cfg.DbQueries.SetSessionIDInvalid(r.Context(), sessionID.ID)
	if err != nil {
		respondWithError(w, r, 500, "could not update database")
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
