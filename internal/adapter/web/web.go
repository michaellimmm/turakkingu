package web

import (
	"context"
	"github/michaellimmm/turakkingu/internal/core"
	"github/michaellimmm/turakkingu/internal/usecase"

	webui "github/michaellimmm/turakkingu/web"
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

var conversionTracker = &webui.ConversionTracker{
	ConversionPoints: []webui.ConversionPoint{
		{ID: "1", Name: "RTG_C99104N99_Shot_P99_ランクル+コンパクトカー訴...", URL: "https://kinto-jp.com/kinto_one/lineup/toyota/?utm_s...", Status: "Draft"},
		{ID: "2", Name: "RTG_C99104N99_Shot_P99_ランクル+コンパクトカー訴...", URL: "https://kinto-jp.com/kinto_one/lineup/toyota/noah/?utm...", Status: "Draft"},
		{ID: "3", Name: "RTG_C99104N99_Shot_P99_ランクル+コンパクトカー訴...", URL: "https://kinto-jp.com/kinto_one/lineup/toyota/noah/?utm...", Status: "Draft"},
		{ID: "4", Name: "RTG_C99104N99_Shot_P99_ランクル+コンパクトカー訴...", URL: "https://kinto-jp.com/kinto_one/lineup/toyota/noah/?utm...", Status: "Draft"},
		{ID: "5", Name: "RTG_C99104N99_Shot_P99_ランクル+コンパクトカー訴...", URL: "https://kinto-jp.com/kinto_one/lineup/toyota/noah/?utm...", Status: "Draft"},
		{ID: "6", Name: "RTG_C99104N99_Shot_P99_ランクル+コンパクトカー訴...", URL: "https://kinto-jp.com/kinto_one/lineup/toyota/noah/?utm...", Status: "Draft"},
	},
}

var landingPageTracker = &webui.LandingPageTracker{
	LandingPages: []webui.LandingPage{
		{ID: "1", FixedURL: "https://example.z...", LandingPageName: "RTG_C99104N99_Shot_P99_ランクル+コ...", LandingPageURL: "https://kinto-jp.com/kinto_one/lineup/toyo..."},
		{ID: "2", FixedURL: "https://example.z...", LandingPageName: "RTG_C99104N99_Shot_P99_ランクル+コ...", LandingPageURL: "https://kinto-jp.com/kinto_one/lineup/toyo..."},
		{ID: "3", FixedURL: "https://example.z...", LandingPageName: "RTG_C99104N99_Shot_P99_ランクル+コ...", LandingPageURL: "https://kinto-jp.com/kinto_one/lineup/toyo..."},
		{ID: "4", FixedURL: "https://example.z...", LandingPageName: "RTG_C99104N99_Shot_P99_ランクル+コ...", LandingPageURL: "https://kinto-jp.com/kinto_one/lineup/toyo..."},
		{ID: "5", FixedURL: "https://example.z...", LandingPageName: "RTG_C99104N99_Shot_P99_ランクル+コ...", LandingPageURL: "https://kinto-jp.com/kinto_one/lineup/toyo..."},
		{ID: "6", FixedURL: "https://example.z...", LandingPageName: "RTG_C99104N99_Shot_P99_ランクル+コ...", LandingPageURL: "https://kinto-jp.com/kinto_one/lineup/toyo..."},
	},
}

type router struct {
	linkWeb         *linkWeb
	thankYouPageWeb *thankYouPageWeb
}

func (r *router) Mux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", r.thankYouPageWeb.Index)
	mux.HandleFunc("POST /search", handleSearch)
	mux.HandleFunc("POST /filter", handleFilter)
	mux.HandleFunc("POST /add", handleAdd)
	mux.HandleFunc("POST /start-tracking", handleStartTracking)

	// Landing pages routes
	mux.HandleFunc("GET /landing-pages", r.linkWeb.Index)
	mux.HandleFunc("POST /landing-pages/search", r.linkWeb.Search)
	mux.HandleFunc("POST /landing-pages/add", r.linkWeb.Create)
	mux.HandleFunc("POST /landing-pages/edit/{id}", r.linkWeb.Edit)
	mux.HandleFunc("POST /landing-pages/bulk-edit", r.linkWeb.BulkEdit)

	return mux
}
