package usecase

import (
	"context"
	"glbackend/internal/entities"
	"glbackend/internal/repo"
	"glbackend/internal/utils"
	"time"
)

type (
	CategoriesUpdateUsecase interface {
		Execute(ctx context.Context, input CategoriesUpdateInput) (entities.Category, error)
	}

	categoriesUpdateUsecase struct {
		repo       repo.Repo
		ctxTimeout time.Duration
	}

	CategoriesUpdateInput struct {
		ID    uint   `json:"-"`
		Name  string `json:"name"`
		Color string `json:"color"`
	}
)

func NewCategoriesUpdateUsecase(
	repo repo.Repo,
	timeout time.Duration,
) CategoriesUpdateUsecase {
	return &categoriesUpdateUsecase{
		repo:       repo,
		ctxTimeout: timeout,
	}
}

func (uc *categoriesUpdateUsecase) Execute(ctx context.Context, input CategoriesUpdateInput) (entities.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	category, err := uc.repo.Grouper().CategoryFindByID(ctx, input.ID)
	if err != nil {
		return entities.Category{}, err
	}

	if err = utils.CopyFields(&input, &category); err != nil {
		return entities.Category{}, err
	}

	return uc.repo.Grouper().CategoryUpdate(ctx, category)
}
