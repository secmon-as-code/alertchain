package lambda

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/m-mizutani/alertchain/pkg/chain"
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

func handleFunctionalURL(base context.Context, chain *chain.Chain, event *events.LambdaFunctionURLRequest) (any, error) {
	s := server.New(func(ctx *model.Context, schema types.Schema, data any) error {
		return chain.HandleAlert(ctx, schema, data)
	})

	w := newHTTPResponseWriter()
	r, err := http.NewRequestWithContext(base, event.RequestContext.HTTP.Method, event.RawPath, nil)
	if err != nil {
		return err.Error(), goerr.Wrap(err, "fail to create http request")
	}

	s.ServeHTTP(w, r)
	if w.code != http.StatusOK {
		return string(w.body), goerr.New("error response from alertchain")
	}

	return string(w.body), nil
}
