package main

import (
	"net/http"
	"time"

	"github.com/sebasukodo/chirpy/internal/auth"
)

type respondNewToken struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {

	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Access Denied")
		return
	}

	refreshToken, err := cfg.dbQueries.GetRefreshTokenByToken(r.Context(), bearer)
	if err != nil {
		respondWithError(w, 401, "Access Denied")
		return
	}

	if refreshToken.RevokedAt.Valid {
		respondWithError(w, 401, "Access Denied")
		return
	}

	if time.Now().UTC().After(refreshToken.ExpiresAt) {
		_, err := cfg.dbQueries.SetRefreshTokenInvalid(r.Context(), refreshToken.Token)
		if err != nil {
			respondWithError(w, 500, "could not update database")
			return
		}
		respondWithError(w, 401, "Access Denied")
		return
	}

	newToken, err := auth.MakeJWT(refreshToken.UserID, cfg.tokenSecret, TokenExpiresInSeconds)
	if err != nil {
		respondWithError(w, 500, "could not create token")
		return
	}

	respondWithJSON(w, 200, respondNewToken{newToken})

}

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {

	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Access Denied")
		return
	}

	refreshToken, err := cfg.dbQueries.GetRefreshTokenByToken(r.Context(), bearer)
	if err != nil {
		respondWithError(w, 401, "Access Denied")
		return
	}

	_, err = cfg.dbQueries.SetRefreshTokenInvalid(r.Context(), refreshToken.Token)
	if err != nil {
		respondWithError(w, 500, "could not update database")
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
