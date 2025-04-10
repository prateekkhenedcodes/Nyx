package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func main() {

	const port = "8080"
	const filePathRoot = "."

	apiCfg := &apiConfig{}

	mux := http.NewServeMux()

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))

	mux.Handle("/app/", apiCfg.MiddleWareMetrics(handler))

	mux.HandleFunc("GET /admin/metrics", apiCfg.CountHits)

	s := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("serving files from %v on port %v", filePathRoot, port)
	s.ListenAndServe()

}
