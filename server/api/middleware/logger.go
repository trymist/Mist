package middleware

import (
	"log"
	"net/http"
)

func Logger(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextHandler.ServeHTTP(w, r)
		log.Printf("=> %s %s", r.Method, r.URL.Path)
	})
}
