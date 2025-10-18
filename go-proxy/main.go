package main

import (
	"go-proxy/config"
	"go-proxy/middleware"
	"io"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("could not load config: %s\n", err)
	}

	proxyHandler := http.HandlerFunc(handleRequest)
	http.Handle("/", middleware.Chain(proxyHandler, middleware.BlockList(cfg.BlockList), middleware.Cookie(), middleware.Cors(), middleware.HeaderRewrite(), middleware.Proxify()))
	log.Println("Proxy server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("could not start server: %s\n", err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Get the target URL from the query string
	targetURL := r.URL.Query().Get("url")
	if targetURL == "" {
		http.Error(w, "url parameter is missing", http.StatusBadRequest)
		return
	}

	// Create a new request to the target URL
	proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy headers from the original request to the new request
	for header, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(header, value)
		}
	}

	// Send the new request
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy headers from the proxy response to the original response writer
	for header, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}

	// Set the status code
	w.WriteHeader(resp.StatusCode)

	// Copy the body from the proxy response to the original response writer
	io.Copy(w, resp.Body)
}