package middleware

import (
	"log"
	"net/http"
	"time"
)

type statusRecorder struct {
	status int
	http.ResponseWriter
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()
		log.Default().Printf("Start of the request: ", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Default().Printf("Completed : %v, %v, %v", r.Method, r.URL.Path, time.Since(start))
	})
}
