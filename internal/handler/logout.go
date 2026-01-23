package handler

import (
	"errors"
	"net/http"

	"github.com/sebasukodo/chirpy/internal/auth"
)

func (cfg *ApiConfig) UserLogout(w http.ResponseWriter, r *http.Request) {

	isRefreshToken := true

	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			isRefreshToken = false
		} else {
			respondWithError(w, r, 400, "Bad Request")
			return
		}
	}

	if isRefreshToken {
		if err := cfg.DbQueries.RevokeRefreshTokenByToken(r.Context(), auth.HashToken(refreshToken.Value)); err != nil {
			respondWithError(w, r, 500, "Internal Error")
			return
		}
	}

	sessionID, err := r.Cookie("session_id")
	if err != nil {
		respondWithError(w, r, 400, "Bad Request")
		return
	}

	if err := cfg.DbQueries.RevokeSessionByID(r.Context(), sessionID.Value); err != nil {
		respondWithError(w, r, 500, "Internal Error")
		return
	}

	cfg.RemoveAllCookies(w)

	w.Header().Set("HX-Redirect", "/login")
	w.WriteHeader(http.StatusOK)

}
