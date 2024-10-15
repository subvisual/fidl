package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func NewZap(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()
			defer func() { //nolint:contextcheck
				status := ww.Status()
				rlog := logger.With(
					zap.String("host", r.RemoteAddr),
					zap.String("proto", r.Proto),
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.String("ua", r.Header.Get("User-Agent")),
					zap.Duration("duration", time.Since(start)),
					zap.String("rid", middleware.GetReqID(r.Context())),
					zap.Int("status", status),
					zap.Int("size", ww.BytesWritten()),
				)

				referer := r.Header.Get("Referer")
				if referer != "" {
					rlog = rlog.With(zap.String("referer", referer))
				}

				query := r.URL.RawQuery
				if query != "" {
					rlog = rlog.With(zap.String("query", query))
				}

				rlog.Info("request")
			}()
			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
