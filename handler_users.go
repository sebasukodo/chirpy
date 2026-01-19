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

type userAuth struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
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

	matching, err := auth.CheckPasswordHash(userLoginRequest.Password, userInfo.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "incorrect email or password")
		return
	}

	if matching {
		respondWithJSON(w, 200, convertDatabaseUser(userInfo))
		return
	}

	respondWithError(w, 401, "incorrect email or password")

}

func convertDatabaseUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
}
