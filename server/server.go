package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := 80
	log.Printf("Starting web server on port %d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), http.FileServer(http.Dir("./web"))); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
