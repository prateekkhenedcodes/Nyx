package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() (string, error) {
	bs := make([]byte, 32)

	_, err := rand.Read(bs)

	if err != nil {
		return "", err
	}

	refreshTokenString := hex.EncodeToString(bs)

	return refreshTokenString, nil
}
