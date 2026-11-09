package main

import (
	"errors"
	"log"
	"net/http"
	"time"
)

// addr is the TCP address the HTTP server listens on.
const addr = ":8080"

func main() {
	srv := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("listening on http://localhost%s", addr)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server stopped: %v", err)
	}
}
