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

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject: userID.String(),
	})

	jwt, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return jwt, nil
}


func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	var claims jwt.RegisteredClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.UUID{}, err
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}

	if subject == "" {
		return uuid.UUID{}, fmt.Errorf("subject claim is empty")
	}

	userID, err := uuid.Parse(subject)
	if err != nil {
		return uuid.UUID{}, err
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	token := headers.Get("Authorization")
	if token == "" {
		return "", fmt.Errorf("authorization header cannot be empty")
	}

	bearerTokenArr := strings.Split(token, " ")
	if len(bearerTokenArr) < 2 || bearerTokenArr[0] != "Bearer" {
		return "", fmt.Errorf("invalid authorization header")
	}

	return strings.TrimSpace(bearerTokenArr[1]), nil
}

func MakeRefreshToken() (string,error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", fmt.Errorf("failed to make refresh token: %w", err)
	}

	return hex.EncodeToString(key), nil
}