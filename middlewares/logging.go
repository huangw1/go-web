package middlewares

import (
	"net/http"
	"time"
	"fmt"
)

func LoggingMiddleware(next http.Handler) http.Handler  {
	n := func(w http.ResponseWriter, r *http.Request) {
		began := time.Now()
		next.ServeHTTP(w, r)
		end := time.Now()
		fmt.Printf("[%s] %q %v", r.Method, r.URL.String(), end.Sub(began))
	}
	return http.HandlerFunc(n)
}