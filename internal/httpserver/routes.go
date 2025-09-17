package httpserver

import (
	"log/slog"
	"net/http"
	"time"
)

func Routes(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	//health endpoint
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r http.Request){
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	return loggingMiddleware(logger, mux)

}

func loggingMiddleware(l *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			l.Info("http_request",
				"method", r.Method,
				"path", r.URL.Path,
				"duration_ms", time.Since(start).Milliseconds(),
			)
	})
}