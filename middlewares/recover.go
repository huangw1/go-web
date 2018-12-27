package middlewares

import (
	"net/http"
	"fmt"
)

func RecoverMiddleware(next http.Handler) http.Handler {
	n := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(n)
}
