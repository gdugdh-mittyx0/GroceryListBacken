package usecase

import (
	"context"
	"glbackend/internal/repo"
	"time"
)

type (
	TagsDeleteUsecase interface {
		Execute(ctx context.Context, params TagsDeleteParams) error
	}

	tagsDeleteUsecase struct {
		repo       repo.Repo
		ctxTimeout time.Duration
	}

	TagsDeleteParams struct {
		ID uint `param:"id"`
	}
)

func NewTagsDeleteUsecase(
	repo repo.Repo,
	timeout time.Duration,
) TagsDeleteUsecase {
	return &tagsDeleteUsecase{
		repo:       repo,
		ctxTimeout: timeout,
	}
}

func (uc *tagsDeleteUsecase) Execute(ctx context.Context, params TagsDeleteParams) error {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	products, err := uc.repo.Product().FindAllByTagID(ctx, params.ID)
	if err != nil {
		return err
	}

	for _, product := range products {
		if err := uc.repo.Grouper().RemoveTagInProduct(ctx, product.ID, params.ID); err != nil {
			return err
		}
	}

	return uc.repo.Grouper().TagDelete(ctx, params.ID)
}
