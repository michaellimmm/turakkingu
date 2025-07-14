package web

import (
	"context"
	"github/michaellimmm/turakkingu/internal/core"
	"github/michaellimmm/turakkingu/internal/entity"
	"github/michaellimmm/turakkingu/internal/usecase"
	"net/http"

	webui "github/michaellimmm/turakkingu/web"

	"github.com/a-h/templ"
)

type linkWeb struct {
	uc     usecase.UseCase
	config *core.Config
}

func NewLinkWeb(config *core.Config, uc usecase.UseCase) *linkWeb {
	return &linkWeb{
		uc:     uc,
		config: config,
	}
}

func (l *linkWeb) Index(w http.ResponseWriter, r *http.Request) {
	var component templ.Component

	// TODO: fix this
	links, err := l.uc.GetAllLinks(r.Context(), "tenant1")
	if err != nil {
		component = webui.LandingPagesContent([]webui.LandingPage{})
	} else {
		res := []webui.LandingPage{}
		for _, link := range links {
			res = append(res, webui.LandingPage{
				ID:              link.ID.Hex(),
				LandingPageName: link.Name,
				LandingPageURL:  link.Url,
				FixedURL:        link.ConstructFixedUrl(l.config.Domain),
			})
		}
		component = webui.LandingPagesContent(res)
	}

	component.Render(context.Background(), w)
}

func (l *linkWeb) Search(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("search")
	var component templ.Component

	links, err := l.uc.SearchLinks(r.Context(), "tenant1", query)
	if err != nil {
		component = webui.LandingPagesTable([]webui.LandingPage{})
	} else {
		res := []webui.LandingPage{}
		for _, link := range links {
			res = append(res, webui.LandingPage{
				ID:              link.ID.Hex(),
				LandingPageName: link.Name,
				LandingPageURL:  link.Url,
				FixedURL:        link.ConstructFixedUrl(l.config.Domain),
			})
		}
		component = webui.LandingPagesTable(res)
	}

	component.Render(context.Background(), w)
}

func (l *linkWeb) Create(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("landing_page_name")
	url := r.FormValue("landing_page_url")

	link := &entity.Link{
		Name:     name,
		Url:      url,
		TenantID: "tenant1",
	}
	_ = l.uc.CreateLink(r.Context(), link)

	var component templ.Component

	// TODO: fix this
	links, err := l.uc.GetAllLinks(r.Context(), "tenant1")
	if err != nil {
		component = webui.LandingPagesTable([]webui.LandingPage{})
	} else {
		res := []webui.LandingPage{}
		for _, link := range links {
			res = append(res, webui.LandingPage{
				ID:              link.ID.Hex(),
				LandingPageName: link.Name,
				LandingPageURL:  link.Url,
				FixedURL:        link.ConstructFixedUrl(l.config.Domain),
			})
		}
		component = webui.LandingPagesTable(res)
	}

	component.Render(context.Background(), w)
}

// TODO: fix this
func (l *linkWeb) Edit(w http.ResponseWriter, r *http.Request) {
	_ = r.PathValue("id")

	_ = r.FormValue("fixed_url")
	_ = r.FormValue("landing_page_name")
	_ = r.FormValue("landing_page_url")

	component := webui.LandingPagesTable([]webui.LandingPage{})
	component.Render(context.Background(), w)
}
