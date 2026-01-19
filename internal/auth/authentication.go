package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	registerClaims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
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

	if claims.Issuer != "chirpy" {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	id := claims.Subject

	uid, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	return uid, nil

}
