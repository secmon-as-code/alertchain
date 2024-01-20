package server

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"

	"log/slog"

	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra/policy"
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
			if authz != nil {
				reader := r.Body
				body, err := io.ReadAll(reader)
				if err != nil {
					utils.HandleError(err)
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
					return
				}
				defer utils.SafeClose(reader)
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
						utils.Logger().Info("rego print",
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
						utils.Logger().Error("Fail to evaluate authz policy", utils.ErrLog(err))
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				}

				if output.Deny {
					w.WriteHeader(http.StatusForbidden)
					utils.SafeWrite(w, []byte("Access denied"))
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
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
