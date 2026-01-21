package handler

import (
	"net/http"
	"time"

	"github.com/sebasukodo/chirpy/internal/auth"
)

const TokenExpiresInSeconds = time.Duration(3600) * time.Second

const RefreshTokenExpiresInHours = time.Duration(60*24) * time.Hour

type respondNewToken struct {
	Token string `json:"token"`
}

func (cfg *ApiConfig) RefreshToken(w http.ResponseWriter, r *http.Request) {

	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Access Denied")
		return
	}

	refreshToken, err := cfg.DbQueries.GetRefreshTokenByToken(r.Context(), bearer)
	if err != nil {
		respondWithError(w, 401, "Access Denied")
		return
	}

	if refreshToken.RevokedAt.Valid {
		respondWithError(w, 401, "Access Denied")
		return
	}

	if time.Now().UTC().After(refreshToken.ExpiresAt) {
		_, err := cfg.DbQueries.SetRefreshTokenInvalid(r.Context(), refreshToken.Token)
		if err != nil {
			respondWithError(w, 500, "could not update database")
			return
		}
		respondWithError(w, 401, "Access Denied")
		return
	}

	newToken, err := auth.MakeJWT(refreshToken.UserID, cfg.TokenSecret, TokenExpiresInSeconds)
	if err != nil {
		respondWithError(w, 500, "could not create token")
		return
	}

	respondWithJSON(w, 200, respondNewToken{newToken})

}

func (cfg *ApiConfig) RevokeToken(w http.ResponseWriter, r *http.Request) {

	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Access Denied")
		return
	}

	refreshToken, err := cfg.DbQueries.GetRefreshTokenByToken(r.Context(), bearer)
	if err != nil {
		respondWithError(w, 401, "Access Denied")
		return
	}

	_, err = cfg.DbQueries.SetRefreshTokenInvalid(r.Context(), refreshToken.Token)
	if err != nil {
		respondWithError(w, 500, "could not update database")
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
