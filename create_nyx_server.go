package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/prateekkhenedcodes/Nyx/internal/auth"
	"github.com/prateekkhenedcodes/Nyx/sql/queries"
)

type NyxServerRes struct {
	ServerId        string `json:"server_id"`
	CreatedAt       string `json:"created_at"`
	ExpiresAt       string `json:"expires_at"`
	MaxParticipants int    `json:"max_participants"`
	ActiveSession   bool   `json:"active_session"`
	UserId          string `json:"user_id"`
}

func (cfg *apiConfig) CreateNyxServer(w http.ResponseWriter, r *http.Request) {

	header := r.Header
	token, err := auth.GetBearerToken(header)
	if err != nil {
		respondWithError(w, 401, "Unauthorised", err)
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.secretToken)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}

	expireTime := 60
	maxParti := 2

	ServerTable, err := queries.AddNyxServer(
		cfg.db,
		uuid.New().String(),
		time.Now().Format(time.RFC3339),
		time.Now().Add(time.Duration(expireTime)*time.Minute).Format(time.RFC3339),
		maxParti,
		true,
		userId)
	if err != nil {
		respondWithError(w, 500, "could not create a nyxServer table", err)
		return
	}

	respondWithJSON(w, 200, NyxServerRes{
		ServerTable.ServerId,
		ServerTable.CreatedAt,
		ServerTable.ExpiresAt,
		ServerTable.MaxParticipants,
		ServerTable.ActiveSession,
		ServerTable.UserId,
	})
}
