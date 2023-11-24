package server

import (
	"net/http"

	"log/slog"

	"github.com/m-mizutani/alertchain/pkg/utils"
)

type StatusCodeWriter struct {
	code int
	http.ResponseWriter
}

func (x *StatusCodeWriter) WriteHeader(code int) {
	x.code = code
	x.ResponseWriter.WriteHeader(code)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := &StatusCodeWriter{ResponseWriter: w}
		next.ServeHTTP(sw, r)
		utils.Logger().Info("request",
			slog.Any("method", r.Method),
			slog.Any("path", r.URL.Path),
			slog.Int("status", sw.code),
			slog.Any("remote", r.RemoteAddr),
		)
	})
}
