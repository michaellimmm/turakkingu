package api

import (
	"context"
	"github/michaellimmm/turakkingu/internal/core"
	"github/michaellimmm/turakkingu/internal/usecase"
	"net/http"
)

type API interface {
	Run() error
	Close(context.Context) error
}

type api struct {
	server *http.Server
}

func NewApi(config *core.Config, uc usecase.UseCase) API {
	linkApi := NewLinkAPI(config, uc)

	router := &router{
		linkApi: linkApi,
	}

	server := &http.Server{
		Addr:    config.HttpPort,
		Handler: router.Mux(),
	}

	return &api{
		server: server,
	}
}

func (a *api) Run() error {
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (a *api) Close(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}

type router struct {
	linkApi LinkAPI
}

func (r *router) Mux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/links", r.linkApi.Create)
	mux.HandleFunc("GET /r/{id}", r.linkApi.Redirect)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.WriteHeader(http.StatusOK)
			return
		}
		render404(w)
	})

	return mux
}
