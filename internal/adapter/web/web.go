package web

import (
	"context"
	"github/michaellimmm/turakkingu/internal/core"
	"github/michaellimmm/turakkingu/internal/usecase"
	"strconv"
	"strings"

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
	router := &router{}
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

var tracker = &webui.ConversionTracker{
	ActiveTab:    "tracking-settings",
	ActiveSubTab: "conversion-point-url",
	ConversionPoints: []webui.ConversionPoint{
		{ID: 1, Name: "RTG_C99104N99_Shot_P99_ランクル+コンパクトカー訴...", URL: "https://kinto-jp.com/kinto_one/lineup/toyota/?utm_s...", Status: "Draft"},
		{ID: 2, Name: "RTG_C99104N99_Shot_P99_ランクル+コンパクトカー訴...", URL: "https://kinto-jp.com/kinto_one/lineup/toyota/noah/?utm...", Status: "Draft"},
		{ID: 3, Name: "RTG_C99104N99_Shot_P99_ランクル+コンパクトカー訴...", URL: "https://kinto-jp.com/kinto_one/lineup/toyota/noah/?utm...", Status: "Draft"},
		{ID: 4, Name: "RTG_C99104N99_Shot_P99_ランクル+コンパクトカー訴...", URL: "https://kinto-jp.com/kinto_one/lineup/toyota/noah/?utm...", Status: "Draft"},
		{ID: 5, Name: "RTG_C99104N99_Shot_P99_ランクル+コンパクトカー訴...", URL: "https://kinto-jp.com/kinto_one/lineup/toyota/noah/?utm...", Status: "Draft"},
		{ID: 6, Name: "RTG_C99104N99_Shot_P99_ランクル+コンパクトカー訴...", URL: "https://kinto-jp.com/kinto_one/lineup/toyota/noah/?utm...", Status: "Draft"},
		{ID: 7, Name: "RTG_C99104N99_Shot_P99_ランクル+コンパクトカー訴...", URL: "https://kinto-jp.com/kinto_one/lineup/toyota/noah/?utm...", Status: "Draft"},
		{ID: 8, Name: "RTG_C99104N99_Shot_P99_ランクル+コンパクトカー訴...", URL: "https://kinto-jp.com/kinto_one/lineup/toyota/noah/?utm...", Status: "Draft"},
		{ID: 9, Name: "RTG_C99104N99_Shot_P99_ランクル+コンパクトカー訴...", URL: "https://kinto-jp.com/kinto_one/lineup/toyota/noah/?utm...", Status: "Draft"},
		{ID: 10, Name: "RTG_C99104N99_Shot_P99_ランクル+コンパクトカー訴...", URL: "https://kinto-jp.com/kinto_one/lineup/toyota/noah/?utm...", Status: "Draft"},
	},
	LandingPages: []webui.LandingPage{
		{ID: 1, FixedURL: "https://example.z...", LandingPageName: "RTG_C99104N99_Shot_P99_ランクル+コ...", LandingPageURL: "https://kinto-jp.com/kinto_one/lineup/toyo...", Status: "Draft"},
		{ID: 2, FixedURL: "https://example.z...", LandingPageName: "RTG_C99104N99_Shot_P99_ランクル+コ...", LandingPageURL: "https://kinto-jp.com/kinto_one/lineup/toyo...", Status: "Draft"},
		{ID: 3, FixedURL: "https://example.z...", LandingPageName: "RTG_C99104N99_Shot_P99_ランクル+コ...", LandingPageURL: "https://kinto-jp.com/kinto_one/lineup/toyo...", Status: "Draft"},
		{ID: 4, FixedURL: "https://example.z...", LandingPageName: "RTG_C99104N99_Shot_P99_ランクル+コ...", LandingPageURL: "https://kinto-jp.com/kinto_one/lineup/toyo...", Status: "Draft"},
		{ID: 5, FixedURL: "https://example.z...", LandingPageName: "RTG_C99104N99_Shot_P99_ランクル+コ...", LandingPageURL: "https://kinto-jp.com/kinto_one/lineup/toyo...", Status: "Draft"},
		{ID: 6, FixedURL: "https://example.z...", LandingPageName: "RTG_C99104N99_Shot_P99_ランクル+コ...", LandingPageURL: "https://kinto-jp.com/kinto_one/lineup/toyo...", Status: "Draft"},
	},
}

type router struct {
}

func (r *router) Mux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handleIndex)
	mux.HandleFunc("GET /tab/{tab}", handleTabSwitch)
	mux.HandleFunc("POST /search", handleSearch)
	mux.HandleFunc("POST /filter", handleFilter)
	mux.HandleFunc("POST /add", handleAdd)
	mux.HandleFunc("POST /start-tracking", handleStartTracking)

	// Landing pages routes
	mux.HandleFunc("POST /landing-pages/search", handleLandingPagesSearch)
	mux.HandleFunc("POST /landing-pages/filter", handleLandingPagesFilter)
	mux.HandleFunc("POST /landing-pages/add", handleLandingPagesAdd)
	mux.HandleFunc("POST /landing-pages/start-tracking", handleLandingPagesStartTracking)

	return mux
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	component := webui.IndexPage(tracker)
	component.Render(context.Background(), w)
}

func handleTabSwitch(w http.ResponseWriter, r *http.Request) {
	tab := r.PathValue("tab")

	tracker.ActiveSubTab = tab

	// Return the entire main content area with updated navigation
	component := webui.MainContentBody(tracker)
	component.Render(context.Background(), w)
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("search")

	filtered := []webui.ConversionPoint{}
	for _, cp := range tracker.ConversionPoints {
		if strings.Contains(strings.ToLower(cp.Name), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(cp.URL), strings.ToLower(query)) {
			filtered = append(filtered, cp)
		}
	}

	component := webui.ConversionPointsTable(filtered)
	component.Render(context.Background(), w)
}

func handleFilter(w http.ResponseWriter, r *http.Request) {
	status := r.FormValue("status")

	filtered := []webui.ConversionPoint{}
	for _, cp := range tracker.ConversionPoints {
		if status == "" || status == "all" || cp.Status == status {
			filtered = append(filtered, cp)
		}
	}

	component := webui.ConversionPointsTable(filtered)
	component.Render(context.Background(), w)
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	url := r.FormValue("url")

	newCP := webui.ConversionPoint{
		ID:     len(tracker.ConversionPoints) + 1,
		Name:   name,
		URL:    url,
		Status: "Draft",
	}

	tracker.ConversionPoints = append(tracker.ConversionPoints, newCP)

	component := webui.ConversionPointsTable(tracker.ConversionPoints)
	component.Render(context.Background(), w)
}

func handleStartTracking(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	selectedIDs := r.Form["selected"]

	for _, idStr := range selectedIDs {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			continue
		}

		for i, cp := range tracker.ConversionPoints {
			if cp.ID == id {
				tracker.ConversionPoints[i].Status = "Active"
				break
			}
		}
	}

	component := webui.ConversionPointsTable(tracker.ConversionPoints)
	component.Render(context.Background(), w)
}

// Landing Pages handlers
func handleLandingPagesSearch(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("search")

	filtered := []webui.LandingPage{}
	for _, lp := range tracker.LandingPages {
		if strings.Contains(strings.ToLower(lp.LandingPageName), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(lp.LandingPageURL), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(lp.FixedURL), strings.ToLower(query)) {
			filtered = append(filtered, lp)
		}
	}

	component := webui.LandingPagesTable(filtered)
	component.Render(context.Background(), w)
}

func handleLandingPagesFilter(w http.ResponseWriter, r *http.Request) {
	status := r.FormValue("status")

	filtered := []webui.LandingPage{}
	for _, lp := range tracker.LandingPages {
		if status == "" || status == "all" || lp.Status == status {
			filtered = append(filtered, lp)
		}
	}

	component := webui.LandingPagesTable(filtered)
	component.Render(context.Background(), w)
}

func handleLandingPagesAdd(w http.ResponseWriter, r *http.Request) {
	fixedURL := r.FormValue("fixed_url")
	landingPageName := r.FormValue("landing_page_name")
	landingPageURL := r.FormValue("landing_page_url")

	newLP := webui.LandingPage{
		ID:              len(tracker.LandingPages) + 1,
		FixedURL:        fixedURL,
		LandingPageName: landingPageName,
		LandingPageURL:  landingPageURL,
		Status:          "Draft",
	}

	tracker.LandingPages = append(tracker.LandingPages, newLP)

	component := webui.LandingPagesTable(tracker.LandingPages)
	component.Render(context.Background(), w)
}

func handleLandingPagesStartTracking(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	selectedIDs := r.Form["selected"]

	for _, idStr := range selectedIDs {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			continue
		}

		for i, lp := range tracker.LandingPages {
			if lp.ID == id {
				tracker.LandingPages[i].Status = "Active"
				break
			}
		}
	}

	component := webui.LandingPagesTable(tracker.LandingPages)
	component.Render(context.Background(), w)
}
