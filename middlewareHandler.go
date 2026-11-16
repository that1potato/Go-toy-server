package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

func middlewareHandler(next http.Handler, timeout time.Duration) http.Handler {
	log.Print("A: middlewareHandler ran.")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeoutCtx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		log.Print("B: inner return in middlewareHandler ran.")
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r.WithContext(timeoutCtx))
	})
}
