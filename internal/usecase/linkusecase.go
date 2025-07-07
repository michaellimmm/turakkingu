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
	GetLink(context.Context, string) (*entity.Link, error)
	GetAllLinks(ctx context.Context, tenantID string) ([]*entity.Link, error)
	SearchLinks(ctx context.Context, tenantID string, keywords string) ([]*entity.Link, error)
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

func (uc *linkUseCase) CreateLink(ctx context.Context, link *entity.Link) error {
	err := uc.repo.CreateLink(ctx, link)
	if err != nil {
		slog.Error("failed to create link", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (uc *linkUseCase) GetLink(ctx context.Context, id string) (*entity.Link, error) {
	link, err := uc.repo.FindLinkByShortID(ctx, id)
	if err != nil {
		slog.Error("failed to get link", slog.String("error", err.Error()))
		return nil, err
	}
	return link, nil
}

func (uc *linkUseCase) GetAllLinks(ctx context.Context, tenantID string) ([]*entity.Link, error) {
	links, err := uc.repo.FindAllLinkbyTenantID(ctx, tenantID)
	if err != nil {
		slog.Error("failed to get links", slog.String("error", err.Error()))
		return nil, err
	}

	return links, err
}

func (uc *linkUseCase) SearchLinks(ctx context.Context, tenantID string, keywords string) ([]*entity.Link, error) {
	links, err := uc.repo.SearchLinks(ctx, tenantID, keywords)
	if err != nil {
		slog.Error("failed to search links", slog.String("error", err.Error()))
		return nil, err
	}

	return links, nil
}
