package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	TokenAccess string = "cirpy-access"
)

func GetBearerToken(headers http.Header) (string, error) {

	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", fmt.Errorf("Access Denied")
	}

	prefix := "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", fmt.Errorf("Acces Denied")
	}

	authHeader = strings.TrimPrefix(authHeader, prefix)

	if authHeader == "" {
		return "", fmt.Errorf("Acces Denied")
	}

	return authHeader, nil

}

func GetAPIKey(headers http.Header) (string, error) {

	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", fmt.Errorf("Access Denied")
	}

	prefix := "ApiKey "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", fmt.Errorf("Access Denied")
	}

	authHeader = strings.TrimPrefix(authHeader, prefix)

	if authHeader == "" {
		return "", fmt.Errorf("Access Denied")
	}

	return authHeader, nil

}

func MakeSessionID() (string, error) {

	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(key), nil

}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	registerClaims := jwt.RegisteredClaims{
		Issuer:    TokenAccess,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, registerClaims)

	jw, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("invalid token")
	}

	return jw, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	claims := jwt.RegisteredClaims{}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	}

	_, err := jwt.ParseWithClaims(tokenString, &claims, keyFunc)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	if claims.Issuer != TokenAccess {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	id := claims.Subject

	uid, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	return uid, nil

}
