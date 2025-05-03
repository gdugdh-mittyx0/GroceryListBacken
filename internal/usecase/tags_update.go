package usecase

import (
	"context"
	"glbackend/internal/entities"
	"glbackend/internal/repo"
	"glbackend/internal/utils"
	"time"
)

type (
	TagsUpdateUsecase interface {
		Execute(ctx context.Context, input TagsUpdateInput) (entities.Tag, error)
	}

	tagsUpdateUsecase struct {
		repo       repo.Repo
		ctxTimeout time.Duration
	}

	TagsUpdateInput struct {
		ID    uint   `json:"-"`
		Name  string `json:"name"`
		Color string `json:"color"`
	}
)

func NewTagsUpdateUsecase(
	repo repo.Repo,
	timeout time.Duration,
) TagsUpdateUsecase {
	return &tagsUpdateUsecase{
		repo:       repo,
		ctxTimeout: timeout,
	}
}

func (uc *tagsUpdateUsecase) Execute(ctx context.Context, input TagsUpdateInput) (entities.Tag, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	tag, err := uc.repo.Grouper().TagFindByID(ctx, input.ID)
	if err != nil {
		return entities.Tag{}, err
	}

	if err = utils.CopyFields(&input, &tag); err != nil {
		return entities.Tag{}, err
	}

	return uc.repo.Grouper().TagUpdate(ctx, tag)
}
