package usecase

import (
	"context"
	"github/michaellimmm/turakkingu/internal/core"
	"github/michaellimmm/turakkingu/internal/entity"
	"github/michaellimmm/turakkingu/internal/repository"
	"log/slog"
)

type LinkUseCase interface {
	CreateLink(context.Context, *entity.Link) error
	GetFixedLink(context.Context, string) (*entity.Link, error)
}

type linkUseCase struct {
	repo   repository.Repo
	config *core.Config
}

func NewLinkUseCase(config *core.Config, repo repository.Repo) LinkUseCase {
	return &linkUseCase{
		repo:   repo,
		config: config,
	}
}

func (l *linkUseCase) CreateLink(ctx context.Context, link *entity.Link) error {
	err := l.repo.Create(ctx, link)
	if err != nil {
		slog.Error("failed to create link", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (l *linkUseCase) GetFixedLink(ctx context.Context, id string) (*entity.Link, error) {
	link, err := l.repo.FindByShortID(ctx, id)
	if err != nil {
		slog.Error("failed to get link", slog.String("error", err.Error()))
		return nil, err
	}
	return link, nil
}
