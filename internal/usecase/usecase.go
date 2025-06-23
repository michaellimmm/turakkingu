package usecase

import (
	"github/michaellimmm/turakkingu/internal/core"
	"github/michaellimmm/turakkingu/internal/repository"
)

type UseCase interface {
	LinkUseCase
}

type usecase struct {
	LinkUseCase
}

func NewUseCase(config *core.Config, repo repository.Repo) UseCase {
	linkUseCase := NewLinkUseCase(config, repo)
	return &usecase{
		LinkUseCase: linkUseCase,
	}
}
