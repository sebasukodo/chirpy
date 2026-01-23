package handler

import (
	"fmt"
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

	if err := cfg.DbQueries.SetSessionIDInvalid(r.Context(), sessionID.ID); err != nil {
		respondWithError(w, r, 401, "Could not delete Session")
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

	if err := cfg.DbQueries.SetSessionIDInvalid(r.Context(), sessionID.ID); err != nil {
		respondWithError(w, r, 401, "Access Denied")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
	})

	w.Header().Set("HX-Redirect", "/login")
	w.WriteHeader(http.StatusOK)

}

func (cfg *ApiConfig) ValidateSessionID(w http.ResponseWriter, r *http.Request) (database.SessionID, error) {

	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		return database.SessionID{}, fmt.Errorf("Access Denied")
	}

	sessionID, err := cfg.DbQueries.GetSessionIDByID(r.Context(), cookie.Value)
	if err != nil {
		return database.SessionID{}, fmt.Errorf("Access Denied")
	}

	if sessionID.RevokedAt.Valid {
		return database.SessionID{}, fmt.Errorf("Session Expired")
	}

	if sessionID.ExpiresAt.Before(time.Now().UTC()) {
		if err := cfg.DbQueries.SetSessionIDInvalid(r.Context(), sessionID.ID); err != nil {
			return database.SessionID{}, fmt.Errorf("Invalidation unsuccessful")
		}
		return database.SessionID{}, fmt.Errorf("Session Expired")
	}

	return sessionID, nil

}

func (cfg *ApiConfig) MakeSession(userInfo database.User, w http.ResponseWriter, r *http.Request) (http.ResponseWriter, database.SessionID, error) {

	sessionID, err := auth.MakeSessionID()
	if err != nil {
		respondWithHTML(templates.LoginError(), w, r)
		return w, database.SessionID{}, err
	}

	session, err := cfg.DbQueries.StoreSessionID(r.Context(), database.StoreSessionIDParams{
		ID:        sessionID,
		UserID:    userInfo.ID,
		ExpiresAt: time.Now().UTC().Add(SessionIDExpiresInHours),
	})
	if err != nil {
		respondWithHTML(templates.LoginError(), w, r)
		return w, database.SessionID{}, err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().UTC().Add(SessionIDExpiresInHours),
	})

	return w, session, nil

}
