package adapter

import (
	"context"
	"github/michaellimmm/turakkingu/internal/adapter/api"
	"github/michaellimmm/turakkingu/internal/adapter/web"
	"github/michaellimmm/turakkingu/internal/core"
	"github/michaellimmm/turakkingu/internal/usecase"

	"golang.org/x/sync/errgroup"
)

type AdapterCloser interface {
	Run() error
	Close(context.Context) error
}

type adapter struct {
	api api.API
	web web.Web
}

func NewAdapter(config *core.Config, uc usecase.UseCase) AdapterCloser {
	api := api.NewApi(config, uc)
	web := web.NewWeb(config, uc)
	return &adapter{
		api: api,
		web: web,
	}
}

func (a *adapter) Run() error {
	var eg errgroup.Group

	eg.Go(a.api.Run)

	eg.Go(a.web.Run)

	return eg.Wait()
}

func (a *adapter) Close(ctx context.Context) error {
	_ = a.api.Close(ctx)
	_ = a.web.Close(ctx)
	return nil
}
