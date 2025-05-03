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
	TagsCreateUsecase interface {
		Execute(ctx context.Context, input TagsCreateInput) (entities.Tag, error)
	}

	tagsCreateUsecase struct {
		repo       repo.Repo
		ctxTimeout time.Duration
	}

	TagsCreateInput struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}
)

func NewTagsCreateUsecase(
	repo repo.Repo,
	timeout time.Duration,
) TagsCreateUsecase {
	return &tagsCreateUsecase{
		repo:       repo,
		ctxTimeout: timeout,
	}
}

func (uc *tagsCreateUsecase) Execute(ctx context.Context, input TagsCreateInput) (entities.Tag, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	if _, err := uc.repo.Grouper().TagFindByName(ctx, input.Name); err != nil {
		if err != gorm.ErrRecordNotFound { // TODO: пред же gorm ErrRecordNotFound использовать. Абстракция ж есть над orm
			return entities.Tag{}, err
		}
	} else {
		return entities.Tag{}, errorsStatus.New(http.StatusBadRequest, "tag exist", "тег существует")
	}

	return uc.repo.Grouper().TagCreate(ctx, entities.Tag{
		Name:  input.Name,
		Color: input.Color,
	})
}
