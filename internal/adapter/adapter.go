package adapter

import (
	"context"
	"github/michaellimmm/turakkingu/internal/adapter/api"
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
}

func NewAdapter(config *core.Config, uc usecase.UseCase) AdapterCloser {
	api := api.NewApi(config, uc)
	return &adapter{
		api: api,
	}
}

func (a *adapter) Run() error {
	var eg errgroup.Group

	eg.Go(a.api.Run)

	return eg.Wait()
}

func (a *adapter) Close(ctx context.Context) error {
	return a.api.Close(ctx)
}
