package usecase

import (
	"context"
	"glbackend/internal/entities"
	"glbackend/internal/repo"
	"time"
)

type (
	ProductsFindAllUsecase interface {
		Execute(ctx context.Context) ([]entities.FullProduct, error)
	}

	productsFindAllUsecase struct {
		repo       repo.Repo
		ctxTimeout time.Duration
	}
)

func NewProductsFindAllUsecase(
	repo repo.Repo,
	timeout time.Duration,
) ProductsFindAllUsecase {
	return &productsFindAllUsecase{
		repo:       repo,
		ctxTimeout: timeout,
	}
}

func (uc *productsFindAllUsecase) Execute(ctx context.Context) ([]entities.FullProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	products, err := uc.repo.Product().FindAll(ctx)
	if err != nil {
		return nil, err
	}

	fullProducts := []entities.FullProduct{}
	for _, product := range products {
		tags, err := uc.repo.Grouper().TagFindAllByProductID(ctx, product.ID)
		if err != nil {
			return nil, err
		}
		category, err := uc.repo.Grouper().CategoryFindByID(ctx, product.CategoryID)
		if err != nil {
			return nil, err
		}
		fullProducts = append(fullProducts, entities.FullProduct{
			Product:  product,
			Tags:     tags,
			Category: category,
		})
	}

	return fullProducts, nil
}
