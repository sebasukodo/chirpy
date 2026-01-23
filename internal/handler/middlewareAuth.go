package handler

import (
	"net/http"
)

func (cfg *ApiConfig) MiddlewareCheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		err := cfg.ValidateAuth(w, r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) MiddlewareCheckAuthLoginPage(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		err := cfg.ValidateAuth(w, r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)

	})
}

func (cfg *ApiConfig) ValidateAuth(w http.ResponseWriter, r *http.Request) error {

	if _, err := cfg.ValidateSessionID(w, r); err == nil {
		return nil
	}

	userID, err := cfg.RotateRefreshToken(w, r)
	if err != nil {
		cfg.RemoveAllCookies(w)
		return err
	}

	_, err = cfg.MakeSession(userID, w, r)
	if err != nil {
		cfg.RemoveAllCookies(w)
		return err
	}

	return nil
}

func (cfg *ApiConfig) RemoveAllCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "refresh_token",
		MaxAge: -1,
		Path:   "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		MaxAge: -1,
		Path:   "/",
	})
}
