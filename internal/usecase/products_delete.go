package usecase

import (
	"context"
	"glbackend/internal/repo"
	"time"
)

type (
	ProductsDeleteUsecase interface {
		Execute(ctx context.Context, params ProductsDeleteParams) error
	}

	productsDeleteUsecase struct {
		repo       repo.Repo
		ctxTimeout time.Duration
	}

	ProductsDeleteParams struct {
		ID uint `param:"id"`
	}
)

func NewProductsDeleteUsecase(
	repo repo.Repo,
	timeout time.Duration,
) ProductsDeleteUsecase {
	return &productsDeleteUsecase{
		repo:       repo,
		ctxTimeout: timeout,
	}
}

func (uc *productsDeleteUsecase) Execute(ctx context.Context, params ProductsDeleteParams) error {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	tags, err := uc.repo.Grouper().TagFindAllByProductID(ctx, params.ID)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		if err := uc.repo.Grouper().RemoveTagInProduct(ctx, params.ID, tag.ID); err != nil {
			return err
		}
	}

	return uc.repo.Product().Delete(ctx, params.ID)
}
