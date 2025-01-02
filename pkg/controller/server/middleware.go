package server

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"

	"log/slog"

	"github.com/secmon-lab/alertchain/pkg/ctxutil"
	"github.com/secmon-lab/alertchain/pkg/domain/interfaces"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
	"github.com/secmon-lab/alertchain/pkg/infra/policy"
	"github.com/secmon-lab/alertchain/pkg/logging"
	"github.com/secmon-lab/alertchain/pkg/utils"
)

type StatusCodeWriter struct {
	code int
	http.ResponseWriter
}

func (x *StatusCodeWriter) WriteHeader(code int) {
	x.code = code
	x.ResponseWriter.WriteHeader(code)
}

type HTTPAuthzInput struct {
	Method string              `json:"method"`
	Path   string              `json:"path"`
	Query  url.Values          `json:"query"`
	Header map[string][]string `json:"header"`
	Remote string              `json:"remote"`
	Body   string              `json:"body"`
	Env    types.EnvVars       `json:"env"`
}

type HTTPAuthzOutput struct {
	Deny bool `json:"deny"`
}

func Authorize(authz *policy.Client, getEnv interfaces.Env) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if authz != nil {
				reader := r.Body
				body, err := io.ReadAll(reader)
				if err != nil {
					utils.HandleError(ctx, err)
					w.WriteHeader(http.StatusBadRequest)
					utils.SafeWrite(ctx, w, []byte(err.Error()))
					return
				}
				defer utils.SafeClose(ctx, reader)
				r.Body = io.NopCloser(bytes.NewReader(body))

				input := &HTTPAuthzInput{
					Method: r.Method,
					Path:   r.URL.Path,
					Query:  r.URL.Query(),
					Header: r.Header,
					Remote: r.RemoteAddr,
					Body:   string(body),
					Env:    getEnv(),
				}

				options := []policy.QueryOption{
					policy.WithPackageSuffix("http"),
					policy.WithRegoPrint(func(file string, row int, msg string) error {
						ctxutil.Logger(ctx).Info("rego print",
							slog.String("file", file),
							slog.Int("row", row),
							slog.String("msg", msg),
							slog.String("package", "authz.http"),
						)
						return nil
					}),
				}

				var output HTTPAuthzOutput
				if err := authz.Query(r.Context(), input, &output, options...); err != nil {
					if !errors.Is(err, types.ErrNoPolicyResult) {
						ctxutil.Logger(ctx).Error("Fail to evaluate authz policy", logging.ErrAttr(err))
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				}

				if output.Deny {
					w.WriteHeader(http.StatusForbidden)
					utils.SafeWrite(ctx, w, []byte("Access denied"))
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := ctxutil.Logger(ctx).With("request_id", types.NewRequestID())
		ctx = ctxutil.InjectLogger(ctx, logger)

		sw := &StatusCodeWriter{ResponseWriter: w}
		next.ServeHTTP(sw, r.WithContext(ctx))
		logger.Info("request",
			slog.Any("method", r.Method),
			slog.Any("path", r.URL.Path),
			slog.Int("status", sw.code),
			slog.Any("remote", r.RemoteAddr),
		)
	})
}
