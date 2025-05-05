package usecase

import (
	"context"
	"glbackend/internal/entities"
	"glbackend/internal/repo"
	"time"
)

type (
	ProductsStatusesUpdateUsecase interface {
		Execute(ctx context.Context, input ProductsStatusesUpdateInput) ([]entities.Product, error)
	}

	productsStatusesUpdateUsecase struct {
		repo       repo.Repo
		ctxTimeout time.Duration
	}

	ProductsStatusesUpdateInput struct {
		Status      entities.StatusProduct `json:"status"`
		ProductsIDs []uint                 `json:"products_ids"`
	}
)

func NewProductsStatusesUpdateUsecase(
	repo repo.Repo,
	timeout time.Duration,
) ProductsStatusesUpdateUsecase {
	return &productsStatusesUpdateUsecase{
		repo:       repo,
		ctxTimeout: timeout,
	}
}

func (uc *productsStatusesUpdateUsecase) Execute(ctx context.Context, input ProductsStatusesUpdateInput) ([]entities.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	products, err := uc.repo.Product().FindAllByIDs(ctx, input.ProductsIDs)
	if err != nil {
		return nil, err
	}

	for i := range products {
		products[i].Status = input.Status
		products[i], err = uc.repo.Product().Update(ctx, products[i])
		if err != nil {
			return nil, err
		}
	}

	return products, nil
}
