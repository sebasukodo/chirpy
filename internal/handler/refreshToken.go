package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sebasukodo/chirpy/internal/auth"
	"github.com/sebasukodo/chirpy/internal/database"
)

const RefreshTokenExpiresInHours = time.Duration(14*24) * time.Hour

func (cfg *ApiConfig) RotateRefreshToken(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return uuid.Nil, nil
		}
		return uuid.Nil, err
	}

	token, err := cfg.DbQueries.GetRefreshTokenByHash(r.Context(), auth.HashToken(cookie.Value))
	if err != nil || token.ExpiresAt.Before(time.Now().UTC()) {
		return uuid.Nil, fmt.Errorf("invalid refresh token")
	}

	err = cfg.DbQueries.SetRefreshTokenInvalid(r.Context(), token.Token)
	if err != nil {
		return uuid.Nil, err
	}

	user, err := cfg.MakeRefreshToken(token.UserID, w, r)
	if err != nil {
		return uuid.Nil, err
	}

	return user, nil
}

func (cfg *ApiConfig) ValidateRefreshToken(w http.ResponseWriter, r *http.Request) (database.RefreshToken, error) {

	cookie, err := r.Cookie("refresh_token")
	if err != nil || cookie.Value == "" {
		return database.RefreshToken{}, fmt.Errorf("Access Denied")
	}

	refreshToken, err := cfg.DbQueries.GetRefreshTokenByHash(r.Context(), auth.HashToken(cookie.Value))
	if err != nil {
		return database.RefreshToken{}, fmt.Errorf("Access Denied")
	}

	if refreshToken.RevokedAt.Valid {
		cfg.DbQueries.RevokeAllRefreshTokensForUser(r.Context(), refreshToken.UserID)
		cfg.DbQueries.RevokeAllSessionsForUser(r.Context(), refreshToken.UserID)
		return database.RefreshToken{}, fmt.Errorf("token reuse detected")
	}

	if refreshToken.ExpiresAt.Before(time.Now().UTC()) {
		if err := cfg.DbQueries.SetRefreshTokenInvalid(r.Context(), refreshToken.Token); err != nil {
			return database.RefreshToken{}, fmt.Errorf("InValidation unsuccessful")
		}
		return database.RefreshToken{}, fmt.Errorf("Token expired")
	}

	return refreshToken, nil

}

func (cfg *ApiConfig) MakeRefreshToken(userID uuid.UUID, w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {

	refreshToken, err := auth.GenerateSecureToken()
	if err != nil {
		return uuid.Nil, err
	}

	hashedToken := auth.HashToken(refreshToken)

	user, err := cfg.DbQueries.StoreRefreshToken(r.Context(), database.StoreRefreshTokenParams{
		HashedToken: hashedToken,
		UserID:      userID,
		ExpiresAt:   time.Now().UTC().Add(RefreshTokenExpiresInHours),
	})
	if err != nil {
		return uuid.Nil, err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().UTC().Add(RefreshTokenExpiresInHours),
	})

	return user.UserID, nil

}
