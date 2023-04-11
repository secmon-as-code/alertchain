package server

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/goerr"
	"golang.org/x/exp/slog"
)

type Server struct {
	mux *chi.Mux
}

func loggingError(msg string, err error) {
	errValues := []any{
		slog.String("error", err.Error()),
	}
	if goErr := goerr.Unwrap(err); goErr != nil {
		for key, value := range goErr.Values() {
			errValues = append(errValues, slog.Any(key, value))
		}
	}

	utils.Logger().Error(msg, errValues...)
}

func respondError(w http.ResponseWriter, err error) {
	body := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}

	loggingError("respond error", err)

	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		utils.Logger().Error("failed to convert error message", err)
		return
	}
}

func getSchema(r *http.Request) (types.Schema, error) {
	schema := chi.URLParam(r, "schema")
	if schema == "" {
		return "", goerr.Wrap(types.ErrInvalidHTTPRequest, "schema is empty")
	}

	return types.Schema(schema), nil
}

func New(route interfaces.Router) *Server {
	s := &Server{}

	handler := func(w http.ResponseWriter, r *http.Request, data any) {
		schema, err := getSchema(r)
		if err != nil {
			respondError(w, err)
			return
		}

		ctx := types.NewContext(types.WithBase(r.Context()))
		if err := route(ctx, schema, data); err != nil {
			respondError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			loggingError("write OK", err)
			return
		}
	}

	r := chi.NewRouter()
	r.Route("/alert", func(r chi.Router) {
		r.Post("/raw/{schema}", func(w http.ResponseWriter, r *http.Request) {
			var data any
			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				respondError(w, err)
				return
			}

			handler(w, r, data)
		})

		r.Post("/pubsub/{schema}", func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				respondError(w, goerr.Wrap(err, "reading pub/sub message").With("body", string(body)))
				return
			}
			utils.Logger().Debug("recv pubsub message", slog.String("body", string(body)))

			var req model.PubSubRequest
			if err := json.Unmarshal(body, &req); err != nil {
				respondError(w, goerr.Wrap(err, "parsing pub/sub message").With("body", string(body)))
				return
			}

			var data any
			if err := json.Unmarshal(req.Message.Data, &data); err != nil {
				respondError(w, goerr.Wrap(err, "parsing pub/sub data field").With("data", string(req.Message.Data)))
				return
			}

			handler(w, r, data)
		})
	})

	s.mux = r

	return s
}

func (x *Server) Run(addr string) error {
	server := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           x.mux,
	}

	if err := server.ListenAndServe(); err != nil {
		return goerr.Wrap(err, "failed to listen")
	}

	return nil
}

func (x *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	x.mux.ServeHTTP(w, r)
}
