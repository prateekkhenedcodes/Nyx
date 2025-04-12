package main

import (
	"fmt"
	"net/http"

	"github.com/prateekkhenedcodes/Nyx/sql/schema"
)

func (cfg *apiConfig) Reset(w http.ResponseWriter, r *http.Request) {

	if cfg.admin != "dev" {
		respondWithError(w, 403, "forbidden", fmt.Errorf("some one tried to reset the database"))
		return
	}

	err := schema.DeleteUserTable(cfg.db)
	if err != nil {
		respondWithError(w, 500, "could not drop the users table", err)
		return
	}
	err = schema.CreateUserTable(cfg.db)
	if err != nil {
		respondWithError(w, 500, "could not create the user table", err)
	}
	cfg.fileServerHits.Store(0)
	respondWithJSON(w, 200, "Reset Successful")
}
