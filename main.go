package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = ":8080"

func main() {

	mux := http.NewServeMux()

	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	server := http.Server{
		Handler: mux,
		Addr:    port,
	}

	fmt.Printf("Server is running on port %v\n", port[1:])
	log.Fatal(server.ListenAndServe())

}
