package middleware

import (
	"bytes"
	"io"
	"marketplace-service/internal/logger"
	"net/http"

	"github.com/sirupsen/logrus"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK, new(bytes.Buffer)}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	lrw.body.Write(b)
	return lrw.ResponseWriter.Write(b)
}

const maxBodySize = 500

func truncateBody(s string) string {
	if len(s) > maxBodySize {
		return s[:maxBodySize] + "..."
	}
	return s
}

func LoggingMiddleware(l logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := l.(*logrus.Logger)

			reqBody, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(reqBody))

			log.WithFields(logrus.Fields{
				"method":  r.Method,
				"path":    r.URL.Path,
				"headers": r.Header,
				"body":    truncateBody(string(reqBody)),
			}).Info("Incoming Request")

			lrw := newLoggingResponseWriter(w)
			next.ServeHTTP(lrw, r)

			log.WithFields(logrus.Fields{
				"method":  r.Method,
				"path":    r.URL.Path,
				"status":  lrw.statusCode,
				"headers": lrw.Header(),
				"body":    truncateBody(lrw.body.String()),
			}).Info("Outgoing Response")
		})
	}
}
