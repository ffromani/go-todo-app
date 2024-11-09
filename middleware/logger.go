package middleware

import (
	"log"
	"net/http"
	"os"
	"time"
)

// Logger logs the API endpoint wall clock execution time
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		logSink := log.New(os.Stdout, "", log.LstdFlags)
		logSink.Printf("middleware: %-6s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
