package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sebasukodo/chirpy/internal/database"
	"github.com/sebasukodo/chirpy/internal/handler"
	"github.com/sebasukodo/chirpy/internal/middleware"
)

const port = "8080"
const filepathRoot = "."

func main() {

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		log.Fatal("DB_URL in .env must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	apiCfg := &handler.ApiConfig{
		FileserverHits: atomic.Int32{},
		DbQueries:      database.New(db),
		Platform:       os.Getenv("PLATFORM"),
		TokenSecret:    os.Getenv("TOKENSECRET"),
		PolkaApiKey:    os.Getenv("POLKA_KEY"),
	}

	mux := http.NewServeMux()
	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", middleware.MetricsInc(apiCfg)(fileServerHandler))

	mux.HandleFunc("GET /api/healthz", handler.Readiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.Metric)

	mux.HandleFunc("POST /api/users", apiCfg.UsersCreate)
	mux.HandleFunc("POST /api/login", apiCfg.UsersLogin)
	mux.HandleFunc("PUT /api/users", apiCfg.UsersChangeCredentials)

	mux.HandleFunc("POST /api/chirps", apiCfg.ChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.ChirpsGetAll)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.ChirpsGetByID)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.ChirpsDeleteByID)

	mux.HandleFunc("POST /admin/reset", apiCfg.Reset)
	mux.HandleFunc("POST /api/refresh", apiCfg.RefreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.RevokeToken)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.VIP)

	server := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	fmt.Printf("Server is running on port %v\n", port)
	log.Fatal(server.ListenAndServe())

}
