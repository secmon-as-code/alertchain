package server

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"log/slog"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi/v5"
	"github.com/m-mizutani/alertchain/pkg/controller/graphql"
	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra/policy"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/goerr"
)

type Server struct {
	mux            *chi.Mux
	authz          *policy.Client
	resolver       *graphql.Resolver
	enableGrappiQL bool
}

type Option func(cfg *Server)

func WithResolver(resolver *graphql.Resolver) Option {
	return func(cfg *Server) {
		cfg.resolver = resolver
	}
}

func WithEnableGraphiQL() Option {
	return func(cfg *Server) {
		cfg.enableGrappiQL = true
	}
}

func WithAuthzPolicy(authz *policy.Client) Option {
	return func(cfg *Server) {
		cfg.authz = authz
	}
}

func respondError(w http.ResponseWriter, err error) {
	body := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}

	sentry.CaptureException(err)

	utils.Logger().Error("respond error", utils.ErrLog(err))

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

func New(route interfaces.Router, options ...Option) *Server {
	s := &Server{}
	for _, opt := range options {
		opt(s)
	}

	wrap := func(handler apiHandler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					respondError(w, goerr.Wrap(err.(error), "panic in handler"))
				}
			}()

			resp, err := handler(r, route)
			if err != nil {
				respondError(w, err)
				return
			}

			body, err := json.Marshal(resp.Data)
			if err != nil {
				respondError(w, err)
				return
			}

			w.WriteHeader(resp.Code)
			if _, err := w.Write(body); err != nil {
				respondError(w, err)
				return
			}
		}
	}

	r := chi.NewRouter()
	r.Use(Logging)
	r.Use(Authorize(s.authz))
	r.Route("/health", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			utils.SafeWrite(w, []byte("OK"))
		})
	})

	r.Route("/alert", func(r chi.Router) {
		r.Post("/raw/{schema}", wrap(handleRawAlert))
		r.Post("/pubsub/{schema}", wrap(handlePubSubAlert))
	})

	if s.resolver != nil {
		gql := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{
			Resolvers: s.resolver,
		}))
		r.Handle("/graphql", gql)

		if s.enableGrappiQL {
			r.Handle("/graphiql", playground.Handler("playground", "/graphql"))
		}
	}

	s.mux = r

	return s
}

type apiResponse struct {
	Code int
	Data any
}

type apiHandler func(r *http.Request, route interfaces.Router) (*apiResponse, error)

func handleRawAlert(r *http.Request, route interfaces.Router) (*apiResponse, error) {
	var data any
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, goerr.Wrap(err, "failed to decode request body")
	}

	schema, err := getSchema(r)
	if err != nil {
		return nil, err
	}

	ctx := model.NewContext(model.WithBase(r.Context()))
	if err := route(ctx, schema, data); err != nil {
		return nil, err
	}

	return &apiResponse{
		Code: http.StatusOK,
		Data: "OK",
	}, nil
}

func handlePubSubAlert(r *http.Request, route interfaces.Router) (*apiResponse, error) {
	schema, err := getSchema(r)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, goerr.Wrap(err, "reading pub/sub message").With("body", string(body))
	}
	utils.Logger().Debug("recv pubsub message", slog.String("body", string(body)))

	var req model.PubSubRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, goerr.Wrap(err, "parsing pub/sub message").With("body", string(body))
	}

	var data any
	if err := json.Unmarshal(req.Message.Data, &data); err != nil {
		return nil, goerr.Wrap(err, "parsing pub/sub data field").With("data", string(req.Message.Data))
	}

	ctx := model.NewContext(model.WithBase(r.Context()))
	if err := route(ctx, schema, data); err != nil {
		return nil, err
	}

	return &apiResponse{
		Code: http.StatusOK,
		Data: "OK",
	}, nil
}

func (x *Server) Run(addr string) error {
	server := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           x.mux,
	}

	defer sentry.Flush(2 * time.Second)

	if err := server.ListenAndServe(); err != nil {
		return goerr.Wrap(err, "failed to listen")
	}

	return nil
}

func (x *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	x.mux.ServeHTTP(w, r)
}
