package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sebasukodo/chirpy/internal/auth"
	"github.com/sebasukodo/chirpy/internal/database"
	"github.com/sebasukodo/chirpy/templates"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
	SessionID   string    `json:"session_id"`
}

type userAuth struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (cfg *ApiConfig) UsersRegisterForm(w http.ResponseWriter, r *http.Request) {

	userInfo := userAuth{
		Password: r.FormValue("password"),
		Email:    r.FormValue("email"),
	}

	hashed, err := auth.HashPassword(userInfo.Password)
	if err != nil {
		respondWithHTML(templates.RegisterError(), w, r)
		return
	}

	userInfoParams := database.CreateUserParams{
		Email:          userInfo.Email,
		HashedPassword: hashed,
	}

	_, err = cfg.DbQueries.CreateUser(r.Context(), userInfoParams)
	if err != nil {
		respondWithHTML(templates.RegisterError(), w, r)
		return
	}

	respondWithHTML(templates.RegisterSuccess(), w, r)

}

func (cfg *ApiConfig) UsersLoginForm(w http.ResponseWriter, r *http.Request) {

	userLoginRequest := userAuth{
		Password: r.FormValue("password"),
		Email:    r.FormValue("email"),
	}

	userInfo, err := cfg.DbQueries.GetUserByEmail(r.Context(), userLoginRequest.Email)
	if err != nil {
		respondWithHTML(templates.LoginError(), w, r)
		return
	}

	check, err := auth.CheckPasswordHash(userLoginRequest.Password, userInfo.HashedPassword)
	if err != nil || !check {
		respondWithHTML(templates.LoginError(), w, r)
		return
	}

	sessionID, err := auth.MakeSessionID()
	if err != nil {
		respondWithHTML(templates.LoginError(), w, r)
		return
	}

	_, err = cfg.DbQueries.StoreSessionID(r.Context(), database.StoreSessionIDParams{
		ID:        sessionID,
		UserID:    userInfo.ID,
		ExpiresAt: time.Now().UTC().Add(SessionIDExpiresInHours),
	})
	if err != nil {
		respondWithHTML(templates.LoginError(), w, r)
		return
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

	user := convertDatabaseUser(userInfo)

	respondWithHTML(templates.LoginSuccess(user.Email), w, r)

}

func (cfg *ApiConfig) UsersChangeCredentials(w http.ResponseWriter, r *http.Request) {

	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, r, 401, "Access Denied")
		return
	}

	decoder := json.NewDecoder(r.Body)

	userRequest := userAuth{}

	if err := decoder.Decode(&userRequest); err != nil {
		respondWithError(w, r, 400, "Bad Request")
		return
	}

	userID, err := auth.ValidateJWT(bearer, cfg.TokenSecret)
	if err != nil {
		respondWithError(w, r, 401, "Access Denied")
		return
	}

	if userRequest.Email != "" {
		if err := cfg.DbQueries.UpdateUserEmail(r.Context(), database.UpdateUserEmailParams{
			ID:    userID,
			Email: userRequest.Email,
		}); err != nil {
			respondWithError(w, r, 500, "incorrect email or password")
			return
		}
	}

	if userRequest.Password != "" {
		hashedPw, err := auth.HashPassword(userRequest.Password)
		if err != nil {
			respondWithError(w, r, 500, "incorrect email or password")
			return
		}

		if err := cfg.DbQueries.UpdateUserPassword(r.Context(), database.UpdateUserPasswordParams{
			ID:             userID,
			HashedPassword: hashedPw,
		}); err != nil {
			respondWithError(w, r, 500, "incorrect email or password")
			return
		}
	}

	userInfo, err := cfg.DbQueries.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, r, 500, "incorrect email or password")
		return
	}

	respondWithJSON(w, 200, convertDatabaseUser(userInfo))

}

func convertDatabaseUser(dbUser database.User) User {
	return User{
		ID:          dbUser.ID,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
		Email:       dbUser.Email,
		IsChirpyRed: dbUser.IsChirpyRed,
	}
}
