package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

const port = ":8080"
const filepathRoot = "."

func main() {

	mux := http.NewServeMux()

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/metrics", apiCfg.handlerMetric)

	mux.HandleFunc("/reset", apiCfg.handlerReset)

	server := http.Server{
		Handler: mux,
		Addr:    port,
	}

	fmt.Printf("Server is running on port %v\n", port[1:])
	log.Fatal(server.ListenAndServe())

}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetric(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)

	line := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	w.Write([]byte(line))

}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)

	cfg.fileserverHits.Swap(0)
	w.Write([]byte("OK"))

}
