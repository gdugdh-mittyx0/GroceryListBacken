package usecase

import (
	"context"
	"glbackend/internal/entities"
	"glbackend/internal/repo"
	"time"
)

type (
	TagsFindAllUsecase interface {
		Execute(ctx context.Context) ([]entities.Tag, error)
	}

	tagsFindAllUsecase struct {
		repo       repo.Repo
		ctxTimeout time.Duration
	}
)

func NewTagsFindAllUsecase(
	repo repo.Repo,
	timeout time.Duration,
) TagsFindAllUsecase {
	return &tagsFindAllUsecase{
		repo:       repo,
		ctxTimeout: timeout,
	}
}

func (uc *tagsFindAllUsecase) Execute(ctx context.Context) ([]entities.Tag, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	return uc.repo.Grouper().TagFindAll(ctx)
}
