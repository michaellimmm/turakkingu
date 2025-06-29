package api

import (
	"context"
	"github/michaellimmm/turakkingu/internal/core"
	"github/michaellimmm/turakkingu/internal/usecase"
	"net/http"

	"github.com/rs/cors"
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
	trackingAPI := NewTrackAPI(config, uc)
	trackingSettingAPI := NewTrackingSettingAPI(config, uc)

	router := &router{
		linkApi:            linkApi,
		trackingAPI:        trackingAPI,
		trackingSettingAPI: trackingSettingAPI,
	}
	server := &http.Server{
		Addr:    config.HttpPort,
		Handler: cors.Default().Handler(router.Mux()),
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
	linkApi            LinkAPI
	trackingAPI        TrackAPI
	trackingSettingAPI TrackingSettingAPI
}

func (r *router) Mux() *http.ServeMux {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("POST /v1/links", r.linkApi.CreateLink)
	mux.HandleFunc("GET /r/{id}", r.linkApi.Redirect)

	mux.HandleFunc("GET /v1/tenants/{tenant_id}/tracking-settings", r.trackingSettingAPI.GetTrackingSetting)
	mux.HandleFunc("POST /v1/tracking-settings/pages", r.trackingSettingAPI.AddThankYouPage)

	mux.HandleFunc("POST /v1/tracks", r.trackingAPI.CreateTrack)
	mux.HandleFunc("POST /v1/tracks/events", r.trackingAPI.TrackEvent)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.WriteHeader(http.StatusOK)
			return
		}
		render404(w)
	})

	return mux
}
