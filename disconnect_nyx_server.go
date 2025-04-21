package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/prateekkhenedcodes/Nyx/internal/auth"
)

func (cfg *apiConfig) DisconnectNyxServer(w http.ResponseWriter, r *http.Request) {

	type parameter struct {
		ServerID string `json:"server_id"`
	}

	header := r.Header
	token, err := auth.GetBearerToken(header)
	if err != nil {
		respondWithError(w, 401, "Unauthorised", err)
		return
	}

	_, err = auth.ValidateJWT(token, cfg.secretToken)
	if err != nil {
		respondWithError(w, 401, "Unauthorised", err)
		return
	}

	params := parameter{}

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "could not decode the parameters", err)
		return
	}

	for conn := range Clients[params.ServerID] {
		if string(conn) == (token) {
			Mu.Lock()
			delete(Clients[params.ServerID], &websocket.Conn{})
			Mu.Unlock()
		}
	}
}
