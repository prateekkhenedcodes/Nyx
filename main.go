package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	mySchema "github.com/prateekkhenedcodes/Nyx/sql/schema"

	_ "github.com/mattn/go-sqlite3"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *sql.DB
	admin          string
	secretToken    string
}

func main() {

	db, err := sql.Open("sqlite3", "./nyx.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	const port = "8080"
	const filePathRoot = "./assets"

	godotenv.Load(".env")
	apiCfg := &apiConfig{}
	apiCfg.db = db
	apiCfg.admin = os.Getenv("ADMIN")
	apiCfg.secretToken = os.Getenv("SECRET_TOKEN")

	err = mySchema.CreateUserTable(apiCfg.db)
	if err != nil {
		log.Fatal(err)
	}
	err = mySchema.CreateRefreshToken(apiCfg.db)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/app/", http.StatusMovedPermanently)
	})
	mux.Handle("/app/", apiCfg.MiddleWareMetrics(handler))
	mux.HandleFunc("GET /admin/metrics", apiCfg.CountHits)
	mux.HandleFunc("GET /api/healthz", ReadinessHandler)
	mux.HandleFunc("POST /api/register", apiCfg.Register)
	mux.HandleFunc("POST /admin/system-reset", apiCfg.Reset)
	mux.HandleFunc("POST /api/login", apiCfg.Login)
	mux.HandleFunc("POST /api/token/refresh", apiCfg.RefreshToken)
	mux.HandleFunc("POST /api/token/revoke", apiCfg.RevokeToken)

	s := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("serving files from %v on port %v", filePathRoot, port)
	log.Fatal(s.ListenAndServe())
}
