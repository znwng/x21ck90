package middleware

import (
	"net/http"
	"time"

	"todo-go/internal/logger"
	"todo-go/internal/status"
	"todo-go/internal/utils"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
	Status int
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.Status = code
	lrw.ResponseWriter.WriteHeader(code)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		status.RequestStarted()
		start := time.Now()

		lrw := &LoggingResponseWriter{
			ResponseWriter: w,
			Status:         200,
		}

		defer func() {
			latency := time.Since(start)
			status.RequestFinished(latency)

			if lrw.Status >= 500 {
				status.RecordErrors()
			}

			logger.Log(
				r.Method,
				lrw.Status,
				r.URL.Path,
				utils.GetClientIP(r),
			)
		}()

		next.ServeHTTP(lrw, r)
	})
}

