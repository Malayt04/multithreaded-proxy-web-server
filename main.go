package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

var cache = NewLRUCache(10)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
targetURL := r.URL

	if targetURL.Scheme == "" || targetURL.Host == "" {
		http.Error(w, "Invalid URL: Missing scheme or host. Please use the proxy with absolute URLs (e.g., curl -x http://localhost:8080 http://example.com)", http.StatusBadRequest)
		log.Printf("Invalid URL request: %s from %s", r.URL.String(), r.RemoteAddr)
		return
	}

	if cachedEntry, found := cache.Get(targetURL.String()); found {
		log.Printf("Cache HIT for: %s %s", r.Method, targetURL.String())
		for name, values := range cachedEntry.Headers {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}
		w.WriteHeader(cachedEntry.StatusCode)
		_, err := w.Write(cachedEntry.Data)
		if err != nil {
			log.Printf("Error writing cached response for %s: %v", targetURL.String(), err)
		}
		return
	}

	log.Printf("Cache MISS for: %s %s. Proxying request...", r.Method, targetURL.String())

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

	var responseBodyBytes []byte
	if resp.Body != nil {
		responseBodyBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body from target %s: %v", targetURL.String(), err)
			http.Error(w, "Failed to read target response body", http.StatusInternalServerError)
			return
		}
	}

	newCacheEntry := CacheEntry{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Data:       responseBodyBytes,
	}

	cache.Put(targetURL.String(), newCacheEntry)
	log.Printf("Cached response for: %s (Status: %d, Size: %d bytes)", targetURL.String(), resp.StatusCode, len(responseBodyBytes))

	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	_, err = w.Write(responseBodyBytes)
	if err != nil {
		log.Printf("Error writing response body to client for %s: %v", targetURL.String(), err)
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