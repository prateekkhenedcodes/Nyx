package main

import (
	"crypto/rand"
	"math/big"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/prateekkhenedcodes/Nyx/sql"
)

type RegisterJSON struct {
	NyxCode string `json:"nyx_code"`
}

func (cfg *apiConfig) Register(w http.ResponseWriter, r *http.Request) {
	const lenghtOfCode = 30
	code, err := generateNyxCode(lenghtOfCode)
	if err != nil {
		respondWithError(w, 500, "Error while generating code", err)
		return
	}

	dbUser, err := sql.AddUser(cfg.db, uuid.New().String(), time.Now().Format(time.RFC3339), time.Now().Format(time.RFC3339), code)
	if err != nil {
		respondWithError(w, 500, "could not insert user data into user table", err)
	}
	respondWithJSON(w, 200, RegisterJSON{
		NyxCode: dbUser.NyxCode,
	})
}

func generateNyxCode(lenght int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, lenght)

	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[n.Int64()]
	}
	return string(code), nil
}
