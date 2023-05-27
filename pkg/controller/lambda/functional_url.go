package lambda

import (
	"bytes"
	"context"
	"encoding/base64"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/m-mizutani/alertchain/pkg/controller/server"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/goerr"
)

type httpResponseWriter struct {
	code   int
	body   []byte
	header map[string][]string
}

func newHTTPResponseWriter() *httpResponseWriter {
	return &httpResponseWriter{
		code:   http.StatusOK,
		header: make(map[string][]string),
	}
}

// Header implements http.ResponseWriter
func (x *httpResponseWriter) Header() http.Header {
	return x.header
}

// Write implements http.ResponseWriter
func (x *httpResponseWriter) Write(d []byte) (int, error) {
	x.body = append(x.body, d...)
	return len(d), nil
}

// WriteHeader implements http.ResponseWriter
func (x *httpResponseWriter) WriteHeader(statusCode int) {
	x.code = statusCode
}

var _ http.ResponseWriter = &httpResponseWriter{}

func NewFunctionalURLHandler() func(ctx context.Context, data any, cb Callback) error {
	return func(ctx context.Context, data any, cb Callback) error {
		var event events.LambdaFunctionURLRequest
		if err := remapEvent(data, &event); err != nil {
			return goerr.Wrap(err, "fail to remap event")
		}
		if event.RawPath == "" {
			return goerr.Wrap(types.ErrInvalidLambdaRequest, "Event is not LambdaFunctionURLRequest")
		}

		s := server.New(func(ctx *model.Context, schema types.Schema, data any) error {
			return cb(ctx, schema, data)
		})

		body, err := base64.StdEncoding.DecodeString(event.Body)
		if err != nil {
			return goerr.Wrap(err, "fail to decode base64")
		}

		w := newHTTPResponseWriter()
		r, err := http.NewRequestWithContext(
			ctx, event.RequestContext.HTTP.Method,
			event.RawPath,
			bytes.NewReader(body),
		)
		if err != nil {
			return goerr.Wrap(err, "fail to create http request")
		}

		s.ServeHTTP(w, r)
		if w.code != http.StatusOK {
			return goerr.New("error response from alertchain")
		}

		return nil
	}
}
