package server

import (
	"encoding/json"
	"net/http"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/go-chi/chi/v5"
	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/goerr"
	"golang.org/x/exp/slog"
)

type Server struct {
	mux *chi.Mux
}

func respondError(w http.ResponseWriter, err error) {
	body := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}

	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		utils.Logger().Error("failed to convert error message", err)
	}
}

func New(route interfaces.Router) *Server {
	s := &Server{}

	handler := func(w http.ResponseWriter, r *http.Request, data any) {
		ctx := types.NewContext(types.WithBase(r.Context()))
		if err := route(ctx, chi.URLParam(r, "label"), data); err != nil {
			utils.Logger().Error("routing data", err, slog.Any("data", data))
			respondError(w, err)
			return
		}
	}

	r := chi.NewRouter()
	r.Route("/alert", func(r chi.Router) {
		r.Post("/raw/{label}", func(w http.ResponseWriter, r *http.Request) {
			var data any
			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				respondError(w, err)
				utils.Logger().Error("parsing alert", err)
				return
			}

			handler(w, r, data)
		})

		r.Post("/pubsub/{label}", func(w http.ResponseWriter, r *http.Request) {
			var msg pubsub.Message
			if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
				respondError(w, goerr.Wrap(err, "parsing pub/sub message"))
				utils.Logger().Error("parsing alert", err)
				return
			}

			var data any
			if err := json.Unmarshal(msg.Data, &data); err != nil {
				respondError(w, goerr.Wrap(err, "parsing pub/sub data field"))
				utils.Logger().Error("parsing alert", err)
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
