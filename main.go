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
const filepathRoot = "./static"

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

	fileServerHandler := http.StripPrefix("/static/", http.FileServer(http.Dir(filepathRoot)))

	mux.Handle("/static/", middleware.CheckAuth(apiCfg)(fileServerHandler))

	mux.HandleFunc("/", apiCfg.Homepage)

	mux.HandleFunc("GET /healthz", handler.Readiness)

	mux.HandleFunc("POST /api/register", apiCfg.UsersRegisterForm)
	mux.HandleFunc("POST /api/login", apiCfg.UsersLoginForm)
	mux.HandleFunc("PUT /api/users", apiCfg.UsersChangeCredentials)

	mux.HandleFunc("GET /register", apiCfg.Register)
	mux.HandleFunc("GET /login", apiCfg.Login)

	mux.HandleFunc("POST /api/chirps", apiCfg.ChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.ChirpsGetAll)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.ChirpsGetByID)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.ChirpsDeleteByID)

	mux.HandleFunc("POST /admin/reset", apiCfg.Reset)
	mux.HandleFunc("POST/api/refresh", apiCfg.RefreshSessionID)
	mux.HandleFunc("POST /api/revoke", apiCfg.RevokeSessionID)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.VIP)

	server := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	fmt.Printf("Server is running on port %v\n", port)
	log.Fatal(server.ListenAndServe())

}
