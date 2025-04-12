package main

import (
	"net/http"
	"time"

	"github.com/prateekkhenedcodes/Nyx/internal/auth"
	"github.com/prateekkhenedcodes/Nyx/sql/queries"
)

func (cfg *apiConfig) RefreshToken(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Token string `json:"token"`
	}

	header := r.Header

	rToken, err := auth.GetBearerToken(header)

	if err != nil {
		respondWithError(w, 401, "Authorization header error", err)
		return
	}

	if rToken == "" {
		respondWithError(w, 401, "Authorization header is empty", err)
		return
	}

	userDt, err := queries.GetUserFromRefreshToken(cfg.db, rToken)
	if err != nil {
		respondWithError(w, 401, "could not find the user id for the refresh token", err)
		return
	}

	if userDt.ExpiresAt == "" {
		respondWithError(w, 401, "the expiration time is empty", err)
		return
	}

	expTime, err := time.Parse(time.RFC3339, userDt.ExpiresAt)
	if err != nil {
		respondWithError(w, 401, "could not parse the srting to covert to time.Time", err)
		return
	}

	if expTime.Before(time.Now()) {
		respondWithError(w, 401, "refresh token is expired", err)
		return
	}

	if userDt.RevokedAt != "" {
		respondWithError(w, 401, "refresh token has expired", err)
		return
	}

	newAccessToken, err := auth.MakeJWT(
		userDt.UserId,
		cfg.secretToken,
		time.Hour)
	if err != nil {
		respondWithError(w, 401, "could not create a new access token", err)
		return
	}

	respondWithJSON(w, 200, response{
		Token: newAccessToken,
	})

}
