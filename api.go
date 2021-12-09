package alertchain

import (
	"net/http"
	"time"

	"github.com/m-mizutani/alertchain/pkg/infra"
	"github.com/m-mizutani/goerr"
)

type apiServer struct {
	clients *infra.Clients
	addr    string
}

func newAPIServer(addr string, clients *infra.Clients) *apiServer {
	return &apiServer{
		clients: clients,
		addr:    addr,
	}
}

func (x *apiServer) Run() error {
	server := &http.Server{
		Addr: x.addr,
		Handler: &apiHandler{
			clients: x.clients,
		},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := server.ListenAndServe(); err != nil {
		return goerr.Wrap(err)
	}

	return nil
}

type apiHandler struct {
	clients *infra.Clients
}

func (x *apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("not found"))
}
