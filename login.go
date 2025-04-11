package main

import (
	"encoding/json"
	"net/http"

	"github.com/prateekkhenedcodes/Nyx/internal/auth"
	"github.com/prateekkhenedcodes/Nyx/sql/queries"
)

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

	respondWithJSON(w, 200, "login successful")
}
