package usecase

import (
	"context"
	"glbackend/internal/entities"
	"glbackend/internal/repo"
	"time"
)

type (
	ProductsFindByIDUsecase interface {
		Execute(ctx context.Context, params ProductsFindByIDParams) (entities.FullProduct, error)
	}

	productsFindByIDUsecase struct {
		repo       repo.Repo
		ctxTimeout time.Duration
	}

	ProductsFindByIDParams struct {
		ID uint `param:"id"`
	}
)

func NewProductsFindByIDUsecase(
	repo repo.Repo,
	timeout time.Duration,
) ProductsFindByIDUsecase {
	return &productsFindByIDUsecase{
		repo:       repo,
		ctxTimeout: timeout,
	}
}

func (uc *productsFindByIDUsecase) Execute(ctx context.Context, params ProductsFindByIDParams) (entities.FullProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	product, err := uc.repo.Product().FindByID(ctx, params.ID)
	if err != nil {
		return entities.FullProduct{}, err
	}

	tags, err := uc.repo.Grouper().TagFindAllByProductID(ctx, product.ID)
	if err != nil {
		return entities.FullProduct{}, err
	}

	category, err := uc.repo.Grouper().CategoryFindByID(ctx, product.CategoryID)
	return entities.FullProduct{
		Product:  product,
		Tags:     tags,
		Category: category,
	}, err
}
