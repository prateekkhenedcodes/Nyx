package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func MakeJWT(id string, secretToken string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "nyx",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   id,
	})
	signedToken, err := token.SignedString([]byte(secretToken))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateJWT(tokenString string, tokenSecret string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return "", fmt.Errorf("unexpected signing method %s", t.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", fmt.Errorf("token is not valid")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}
	userId, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("claim does not contain user ID")
	}
	return userId, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	token, err := stripBearer(auth)
	if err != nil {
		return "", err
	}
	return token, nil

}

func stripBearer(s string) (string, error) {
	token := strings.TrimSpace(strings.TrimPrefix(s, "Bearer "))
	if token == "" {
		return "", fmt.Errorf("token string deosn't exist")
	}
	return token, nil
}
