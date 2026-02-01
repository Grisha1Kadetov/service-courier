package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/metrics"
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/pkg/log"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func NewMetricsMiddleware(l log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/metrics" {
					next.ServeHTTP(w, r)
					return
				}
				rw := newLoggingResponseWriter(w)
				start := time.Now()
				next.ServeHTTP(rw, r)
				duration := time.Since(start)
				metrics.RequestsTotal.WithLabelValues(
					r.URL.Path,
					r.Method,
					strconv.Itoa(rw.statusCode),
				).Inc()
				metrics.RequestDuration.WithLabelValues(
					r.URL.Path,
					r.Method,
					strconv.Itoa(rw.statusCode),
				).Observe(duration.Seconds())
				l.Info("", log.NewField("method", r.Method), log.NewField("path", r.URL.Path), log.NewField("status", rw.statusCode), log.NewField("duration", duration))
			})
	}
}
