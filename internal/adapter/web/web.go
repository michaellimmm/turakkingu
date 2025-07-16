package web

import (
	"context"
	"github/michaellimmm/turakkingu/internal/core"
	"github/michaellimmm/turakkingu/internal/usecase"

	"net/http"
)

type Web interface {
	Run() error
	Close(context.Context) error
}

type web struct {
	server *http.Server
}

func NewWeb(config *core.Config, uc usecase.UseCase) Web {
	linkWeb := NewLinkWeb(config, uc)
	thankYouPageWeb := NewThankYouPageWeb(config, uc)
	router := &router{linkWeb: linkWeb, thankYouPageWeb: thankYouPageWeb}
	server := &http.Server{
		Addr:    config.WebPort,
		Handler: router.Mux(),
	}
	return &web{
		server: server,
	}
}

func (w *web) Run() error {
	if err := w.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (w *web) Close(ctx context.Context) error {
	return w.server.Shutdown(ctx)
}

type router struct {
	linkWeb         *linkWeb
	thankYouPageWeb *thankYouPageWeb
}

func (r *router) Mux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", r.thankYouPageWeb.Index)
	mux.HandleFunc("POST /search", handleSearch)
	mux.HandleFunc("POST /filter", r.thankYouPageWeb.handleFilter)
	mux.HandleFunc("POST /add", handleAdd)
	mux.HandleFunc("POST /start-tracking", handleStartTracking)

	// Landing pages routes
	mux.HandleFunc("GET /landing-pages", r.linkWeb.Index)
	mux.HandleFunc("POST /landing-pages/search", r.linkWeb.Search)
	mux.HandleFunc("POST /landing-pages/add", r.linkWeb.Create)
	mux.HandleFunc("POST /landing-pages/edit/{id}", r.linkWeb.Edit)

	return mux
}
