package main

import (
	"net/http"
	"time"

	"github.com/prateekkhenedcodes/Nyx/internal/auth"
	"github.com/prateekkhenedcodes/Nyx/sql/queries"
)

func (cfg *apiConfig) RevokeToken(w http.ResponseWriter, r *http.Request) {
	header := r.Header

	rToken, err := auth.GetBearerToken(header)
	if err != nil {
		respondWithError(w, 401, "Authorization header is empty", err)
		return
	}

	_, err = queries.RevokeRefreshToken(cfg.db,
		rToken,
		time.Now().Format(time.RFC3339),
		time.Now().Format(time.RFC3339))
	if err != nil {
		respondWithError(w, 401, "could not find any user of the token", err)
		return
	}

	respondWithJSON(w, 204, nil)

}
