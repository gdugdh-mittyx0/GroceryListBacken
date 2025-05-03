package usecase

import (
	"context"
	"glbackend/internal/entities"
	"glbackend/internal/errorsStatus"
	"glbackend/internal/repo"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type (
	CategoriesCreateUsecase interface {
		Execute(ctx context.Context, input CategoriesCreateInput) (entities.Category, error)
	}

	categoriesCreateUsecase struct {
		repo       repo.Repo
		ctxTimeout time.Duration
	}

	CategoriesCreateInput struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}
)

func NewCategoriesCreateUsecase(
	repo repo.Repo,
	timeout time.Duration,
) CategoriesCreateUsecase {
	return &categoriesCreateUsecase{
		repo:       repo,
		ctxTimeout: timeout,
	}
}

func (uc *categoriesCreateUsecase) Execute(ctx context.Context, input CategoriesCreateInput) (entities.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	if _, err := uc.repo.Grouper().CategoryFindByName(ctx, input.Name); err != nil {
		if err != gorm.ErrRecordNotFound { // TODO: пред же gorm ErrRecordNotFound использовать. Абстракция ж есть над orm
			return entities.Category{}, err
		}
	} else {
		return entities.Category{}, errorsStatus.New(http.StatusBadRequest, "category exist", "тег существует")
	}

	return uc.repo.Grouper().CategoryCreate(ctx, entities.Category{
		Name:  input.Name,
		Color: input.Color,
	})
}
