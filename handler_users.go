package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sebasukodo/chirpy/internal/auth"
	"github.com/sebasukodo/chirpy/internal/database"
)

const TokenExpiresInSeconds = time.Duration(3600) * time.Second
const RefreshTokenExpiresInHours = time.Duration(60*24) * time.Hour

type userAuth struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	userInfo := userAuth{}

	if err := decoder.Decode(&userInfo); err != nil {
		respondWithError(w, 500, fmt.Sprintf("error while decoding: %v", err))
		return
	}

	hashed, err := auth.HashPassword(userInfo.Password)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("could not hash password: %v", err))
		return
	}

	userInfoParams := database.CreateUserParams{
		Email:          userInfo.Email,
		HashedPassword: hashed,
	}

	dbUser, err := cfg.dbQueries.CreateUser(r.Context(), userInfoParams)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("could not create user: %v", err))
		return
	}

	respondWithJSON(w, 201, convertDatabaseUser(dbUser))

}

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	userLoginRequest := userAuth{}

	if err := decoder.Decode(&userLoginRequest); err != nil {
		respondWithError(w, 401, "incorrect email or password")
		return
	}

	userInfo, err := cfg.dbQueries.GetUserByEmail(r.Context(), userLoginRequest.Email)
	if err != nil {
		respondWithError(w, 401, "incorrect email or password")
		return
	}

	check, err := auth.CheckPasswordHash(userLoginRequest.Password, userInfo.HashedPassword)
	if err != nil || !check {
		respondWithError(w, 401, "incorrect email or password")
		return
	}

	jwtToken, err := auth.MakeJWT(userInfo.ID, cfg.tokenSecret, TokenExpiresInSeconds)
	if err != nil {
		respondWithError(w, 401, "incorrect email or password")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 401, "incorrect email or password")
		return
	}

	_, err = cfg.dbQueries.StoreRefreshToken(r.Context(), database.StoreRefreshTokenParams{
		Token:     refreshToken,
		UserID:    userInfo.ID,
		ExpiresAt: time.Now().UTC().Add(RefreshTokenExpiresInHours),
	})
	if err != nil {
		respondWithError(w, 500, "could not store refresh token")
		return
	}

	userWithToken := convertDatabaseUser(userInfo)
	userWithToken.Token = jwtToken
	userWithToken.RefreshToken = refreshToken

	respondWithJSON(w, 200, userWithToken)

}

func (cfg *apiConfig) handlerUsersChangeCredentials(w http.ResponseWriter, r *http.Request) {

	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Access Denied")
		return
	}

	decoder := json.NewDecoder(r.Body)

	userRequest := userAuth{}

	if err := decoder.Decode(&userRequest); err != nil {
		respondWithError(w, 400, "Bad Request")
		return
	}

	userID, err := auth.ValidateJWT(bearer, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, 401, "Access Denied")
		return
	}

	if userRequest.Email != "" {
		if err := cfg.dbQueries.UpdateUserEmail(r.Context(), database.UpdateUserEmailParams{
			ID:    userID,
			Email: userRequest.Email,
		}); err != nil {
			respondWithError(w, 500, "incorrect email or password")
			return
		}
	}

	if userRequest.Password != "" {
		hashedPw, err := auth.HashPassword(userRequest.Password)
		if err != nil {
			respondWithError(w, 500, "incorrect email or password")
			return
		}

		if err := cfg.dbQueries.UpdateUserPassword(r.Context(), database.UpdateUserPasswordParams{
			ID:             userID,
			HashedPassword: hashedPw,
		}); err != nil {
			respondWithError(w, 500, "incorrect email or password")
			return
		}
	}

	userInfo, err := cfg.dbQueries.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, 500, "incorrect email or password")
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
