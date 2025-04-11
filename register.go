package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/prateekkhenedcodes/Nyx/internal/auth"
	"github.com/prateekkhenedcodes/Nyx/sql/queries"
)

type RegisterJSON struct {
	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	NyxCode   string `json:"nyx_code"`
}

func (cfg *apiConfig) Register(w http.ResponseWriter, r *http.Request) {
	const lenghtOfCode = 30
	code, err := auth.GenerateNyxCode(lenghtOfCode)
	if err != nil {
		respondWithError(w, 500, "Error while generating code", err)
		return
	}

	hashedCode, err := auth.HashNyxCode(code)
	if err != nil {
		respondWithError(w, 500, "could not hash the Nyx code", err)
		return
	}

	dbUser, err := queries.AddUser(cfg.db, uuid.New().String(), time.Now().Format(time.RFC3339), time.Now().Format(time.RFC3339), hashedCode)
	if err != nil {
		respondWithError(w, 500, "could not insert user data into user table", err)
		return
	}
	respondWithJSON(w, 200, RegisterJSON{
		Id:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		NyxCode:   code,
	})
}
