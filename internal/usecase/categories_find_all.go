package usecase

import (
	"context"
	"glbackend/internal/entities"
	"glbackend/internal/repo"
	"time"
)

type (
	CategoriesFindAllUsecase interface {
		Execute(ctx context.Context) ([]entities.Category, error)
	}

	categoriesFindAllUsecase struct {
		repo       repo.Repo
		ctxTimeout time.Duration
	}
)

func NewCategoriesFindAllUsecase(
	repo repo.Repo,
	timeout time.Duration,
) CategoriesFindAllUsecase {
	return &categoriesFindAllUsecase{
		repo:       repo,
		ctxTimeout: timeout,
	}
}

func (uc *categoriesFindAllUsecase) Execute(ctx context.Context) ([]entities.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	return uc.repo.Grouper().CategoryFindAll(ctx)
}
