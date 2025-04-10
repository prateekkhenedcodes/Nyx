package main

import "net/http"

type hits struct {
	TotalHits int `json:"total_hits"`
}

func (cfg *apiConfig) MiddleWareMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) CountHits(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, hits{
		TotalHits: int(cfg.fileServerHits.Load()),
	})
}
