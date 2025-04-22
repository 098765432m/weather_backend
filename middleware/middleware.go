package middleware

import (
	"log"
	"net/http"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request URI: %s\n", r.RequestURI)

		next.ServeHTTP(w, r)
	})
}
