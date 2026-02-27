package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/play/texas-holdem", corsMiddleware(texasHoldem))

	log.Printf("Games API Server starting on :%s", port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

}
