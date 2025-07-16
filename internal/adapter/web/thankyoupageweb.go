package web

import (
	"context"
	"github/michaellimmm/turakkingu/internal/core"
	"github/michaellimmm/turakkingu/internal/usecase"
	webui "github/michaellimmm/turakkingu/web"
	"net/http"
)

type thankYouPageWeb struct {
	uc     usecase.UseCase
	config *core.Config
}

func NewThankYouPageWeb(config *core.Config, uc usecase.UseCase) *thankYouPageWeb {
	return &thankYouPageWeb{
		uc:     uc,
		config: config,
	}
}

// TODO: fix this
func (t *thankYouPageWeb) Index(w http.ResponseWriter, r *http.Request) {
	trackingSetting, _ := t.uc.GetTrackingSettingByTenantID(r.Context(), "tenant1")

	res := &webui.ConversionTracker{}
	res.ConversionPoints = []webui.ConversionPoint{}
	for _, page := range trackingSetting.ThankYouPages {
		res.ConversionPoints = append(res.ConversionPoints, webui.ConversionPoint{
			ID:     page.ID.Hex(),
			Name:   page.Name,
			URL:    page.URL,
			Status: "Draft",
		})
	}

	component := webui.IndexPage(res)
	component.Render(context.Background(), w)
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	_ = r.FormValue("search")

	component := webui.ConversionPointsTable([]webui.ConversionPoint{})
	component.Render(context.Background(), w)
}

func (t *thankYouPageWeb) handleFilter(w http.ResponseWriter, r *http.Request) {
	_ = r.FormValue("status")
	tenantID := r.FormValue("tenant_id")

	trackingSetting, _ := t.uc.GetTrackingSettingByTenantID(r.Context(), tenantID)

	res := &webui.ConversionTracker{}
	res.ConversionPoints = []webui.ConversionPoint{}
	for _, page := range trackingSetting.ThankYouPages {
		res.ConversionPoints = append(res.ConversionPoints, webui.ConversionPoint{
			ID:     page.ID.Hex(),
			Name:   page.Name,
			URL:    page.URL,
			Status: "Draft",
		})
	}

	component := webui.ConversionPointsTable(res.ConversionPoints)
	component.Render(context.Background(), w)
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	_ = r.FormValue("name")
	_ = r.FormValue("url")

	component := webui.ConversionPointsTable([]webui.ConversionPoint{})
	component.Render(context.Background(), w)
}

func handleStartTracking(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	_ = r.Form["selected"]

	component := webui.ConversionPointsTable([]webui.ConversionPoint{})
	component.Render(context.Background(), w)
}
