package server

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"log/slog"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
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
	env            interfaces.Env
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

func WithEnv(env interfaces.Env) Option {
	return func(cfg *Server) {
		cfg.env = env
	}
}

func respondError(w http.ResponseWriter, err error) {
	body := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}

	utils.HandleError(err)

	var code int
	switch types.GetErrorType(err) {
	case
		types.ErrTypeAction,
		types.ErrTypePolicy,
		types.ErrTypeRuntime,
		types.ErrTypeConfig,
		types.ErrTypeUnknown:
		code = http.StatusInternalServerError

	case types.ErrTypeBadRequest:
		code = http.StatusBadRequest
	}

	w.WriteHeader(code)
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

func New(hdlr interfaces.AlertHandler, options ...Option) *Server {
	s := &Server{
		env: utils.Env,
	}
	for _, opt := range options {
		opt(s)
	}

	wrap := func(handler apiAlertHandler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					respondError(w, goerr.Wrap(err.(error), "panic in handler"))
				}
			}()

			resp, err := handler(r, hdlr)
			if err != nil {
				respondError(w, err)
				return
			}

			body := struct {
				Alerts []*model.Alert `json:"alerts"`
			}{
				Alerts: resp.Alerts,
			}

			w.WriteHeader(resp.Code)
			if err := json.NewEncoder(w).Encode(body); err != nil {
				respondError(w, goerr.Wrap(err, "failed to encode response body"))
				return
			}
		}
	}

	r := chi.NewRouter()
	r.Use(Logging)
	r.Use(Authorize(s.authz, s.env))
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

type apiAlertResponse struct {
	Code   int
	Alerts []*model.Alert
}

type apiAlertHandler func(r *http.Request, route interfaces.AlertHandler) (*apiAlertResponse, error)

func handleRawAlert(r *http.Request, route interfaces.AlertHandler) (*apiAlertResponse, error) {
	var data any
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, goerr.Wrap(err, "failed to decode request body")
	}

	schema, err := getSchema(r)
	if err != nil {
		return nil, err
	}

	ctx := model.NewContext(model.WithBase(r.Context()))
	alerts, err := route(ctx, schema, data)
	if err != nil {
		return nil, err
	}

	return &apiAlertResponse{
		Code:   http.StatusOK,
		Alerts: alerts,
	}, nil
}

func handlePubSubAlert(r *http.Request, route interfaces.AlertHandler) (*apiAlertResponse, error) {
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
	alerts, err := route(ctx, schema, data)
	if err != nil {
		return nil, err
	}

	return &apiAlertResponse{
		Code:   http.StatusOK,
		Alerts: alerts,
	}, nil
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
