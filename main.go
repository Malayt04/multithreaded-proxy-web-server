package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	targetURL := r.URL

	log.Printf("Proxying request for: %s %s from %s", r.Method, targetURL.String(), r.RemoteAddr)

	newRequest, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
	if err != nil {
		http.Error(w, "Failed to create new request", http.StatusInternalServerError)
		log.Printf("Error creating new request: %v", err)
		return
	}

	for name, values := range r.Header {
		if name == "Connection" || name == "Proxy-Authenticate" ||
		   name == "Proxy-Authorization" || name == "Te" ||
		   name == "Trailers" || name == "Transfer-Encoding" ||
		   name == "Upgrade" {
			continue
		}
		for _, value := range values {
			newRequest.Header.Add(name, value)
		}
	}

	resp, err := http.DefaultClient.Do(newRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error connecting to target server: %v", err), http.StatusBadGateway)
		log.Printf("Error performing request to target (%s): %v", targetURL.String(), err)
		return
	}
	defer resp.Body.Close()

	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	w.WriteHeader(resp.StatusCode)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Error copying response body for %s: %v", targetURL.String(), err)
	}

	log.Printf("Proxied %s %s with status %d", r.Method, targetURL.String(), resp.StatusCode)
}

func main() {
	http.HandleFunc("/", proxyHandler)

	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}