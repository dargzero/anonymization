package swagger

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

func logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

func safetyNet(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Error catched by global error handler: %v\n%s", r, debug.Stack())
				respondWithError(w, http.StatusInternalServerError, fmt.Sprint(r))
			}
		}()

		inner.ServeHTTP(w, r)
	})
}
