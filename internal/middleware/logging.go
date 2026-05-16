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
		// 1. Start the stopwatch
		start := time.Now()

		// 2. Wrap the original writer in our custom recorder
		// We default to 200 OK, because if WriteHeader is never explicitly called, Go defaults to 200.
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		// 3. Pass the request to your actual application logic
		next.ServeHTTP(recorder, r)

		// 4. Stop the stopwatch and log the results
		log.Printf("%s %s %d %v", r.Method, r.URL.Path, recorder.status, time.Since(start))
	})
}
