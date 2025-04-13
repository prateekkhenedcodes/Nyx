package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/prateekkhenedcodes/Nyx/internal/auth"
	"github.com/prateekkhenedcodes/Nyx/sql/queries"
)

type Login struct {
	ID           string `json:"id"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	JWToken      string `json:"jwt_token"`
	RefreshTOken string `json:"refresh_token"`
}

func (cfg *apiConfig) Login(w http.ResponseWriter, r *http.Request) {
	type loginPara struct {
		Id      string `json:"id"`
		NyxCode string `json:"nyx_code"`
	}

	decoder := json.NewDecoder(r.Body)
	params := loginPara{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 500, "could not decode the request", err)
		return
	}

	dbUser, err := queries.GetUserById(cfg.db, params.Id)
	if err != nil {
		respondWithError(w, 404, "User id doesn't exists, check id", err)
		return
	}

	err = auth.CheckPassHash(dbUser.NyxCode, params.NyxCode)
	if err != nil {
		respondWithError(w, 401, "invalid email or password", err)
		return
	}

	refreshTokenData, err := queries.GetRefreshTokenbyId(cfg.db, dbUser.ID)
	if err != nil {
		respondWithError(w, 500, "could not get the refrestoken by id ", err)
		return
	}

	if refreshTokenData.RevokedAt == "" {
		respondWithError(w, 409, "you are already logged in", fmt.Errorf("your refresh token is not revoked"))
		return
	}

	expTime, err := time.Parse(time.RFC3339, refreshTokenData.ExpiresAt)
	if err != nil {
		respondWithError(w, 500, "could not parse the time string to formatted type", err)
		return
	}

	if !expTime.Before(time.Now()) {
		respondWithError(w, 409, "you are already logged in ", fmt.Errorf("your refresh token is still valid"))
		return
	}

	defaultExpiryTime := 3600

	token, err := auth.MakeJWT(dbUser.ID, cfg.secretToken, time.Duration(defaultExpiryTime)*time.Second)
	if err != nil {
		respondWithError(w, 500, "could not make a JWT ", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 500, "could not generate a new refresh token", err)
	}

	refreshTData, err := queries.AddRefreshToken(cfg.db,
		refreshToken,
		time.Now().Format(time.RFC3339),
		time.Now().Format(time.RFC3339),
		dbUser.ID,
		time.Now().AddDate(0, 0, 60).Format(time.RFC3339),
		"")
	if err != nil {
		respondWithError(w, 500, "could not add refresh token to data table", err)
		return
	}

	respondWithJSON(w, 201, Login{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		JWToken:      token,
		RefreshTOken: refreshTData.Token,
	})
}
