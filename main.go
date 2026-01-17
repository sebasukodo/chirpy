package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = ":8080"

func main() {

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	mux.Handle("/asstes/", http.FileServer(http.Dir("./assets")))

	server := http.Server{
		Handler: mux,
		Addr:    port,
	}

	fmt.Printf("Starting Server on Port %v\n", port[1:])
	log.Fatal(server.ListenAndServe())

}
