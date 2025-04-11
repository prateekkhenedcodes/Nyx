package main

import (
	"database/sql"
	"log"
	"net/http"
	"sync/atomic"

	_ "github.com/mattn/go-sqlite3"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db *sql.DB
}

func main() {

	db, err := sql.Open("sqlite3", "./nyx.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	const port = "8080"
	const filePathRoot = "."

	apiCfg := &apiConfig{}
	apiCfg.db = db

	mux := http.NewServeMux()

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))

	mux.Handle("/app/", apiCfg.MiddleWareMetrics(handler))
	mux.HandleFunc("GET /admin/metrics", apiCfg.CountHits)
	mux.HandleFunc("GET /api/healthz", ReadinessHandler)
	mux.HandleFunc("POST /api/register", apiCfg.Register)

	s := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("serving files from %v on port %v", filePathRoot, port)
	s.ListenAndServe()

}
