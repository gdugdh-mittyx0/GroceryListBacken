package usecase

import (
	"context"
	"glbackend/internal/repo"
	"time"
)

type (
	CategoriesDeleteUsecase interface {
		Execute(ctx context.Context, params CategoriesDeleteParams) error
	}

	categoriesDeleteUsecase struct {
		repo       repo.Repo
		ctxTimeout time.Duration
	}

	CategoriesDeleteParams struct {
		ID uint `param:"id"`
	}
)

func NewCategoriesDeleteUsecase(
	repo repo.Repo,
	timeout time.Duration,
) CategoriesDeleteUsecase {
	return &categoriesDeleteUsecase{
		repo:       repo,
		ctxTimeout: timeout,
	}
}

func (uc *categoriesDeleteUsecase) Execute(ctx context.Context, params CategoriesDeleteParams) error {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	return uc.repo.Grouper().CategoryDelete(ctx, params.ID)
}
