package main

import (
	"fmt"
	"log"
	"net/http"
)

func handleHello(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Hello, client! You requested: %s\n", r.URL.Path)
	log.Printf("Received request for: %s from %s", r.URL.Path, r.RemoteAddr)
}
func main() {

	http.HandleFunc("/", handleHello)

	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}