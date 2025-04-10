package main

import (
	"net/http"
)

func ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, http.StatusOK)
}
