package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashNyxCode(code string) (string, error) {

	const highCost = 12

	hash, err := bcrypt.GenerateFromPassword([]byte(code), highCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckPassHash(hash string, code string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(code))
	if err != nil {
		return fmt.Errorf("passwords dont match: %s", err)
	}
	return nil
}
