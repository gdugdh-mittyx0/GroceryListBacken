package usecase

import (
	"context"
	"glbackend/internal/entities"
	"glbackend/internal/errorsStatus"
	"glbackend/internal/repo"
	"net/http"
	"time"
)

type (
	ProductsCreateUsecase interface {
		Execute(ctx context.Context, input ProductsCreateInput) (entities.FullProduct, error)
	}

	productsCreateUsecase struct {
		repo       repo.Repo
		ctxTimeout time.Duration
	}

	ProductsCreateInput struct {
		Name       string `json:"name"`
		Priority   int    `json:"priority"`
		Icon       string `json:"icon"`
		CategoryID uint   `json:"category_id"`
		TagsID     []uint `json:"tags_id"`
	}
)

func NewProductsCreateUsecase(
	repo repo.Repo,
	timeout time.Duration,
) ProductsCreateUsecase {
	return &productsCreateUsecase{
		repo:       repo,
		ctxTimeout: timeout,
	}
}

func (uc *productsCreateUsecase) Execute(ctx context.Context, input ProductsCreateInput) (entities.FullProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	// TODO: add transaction
	tags, err := uc.repo.Grouper().TagFindAll(ctx)
	if err != nil {
		return entities.FullProduct{}, err
	}

	tagsMap := make(map[uint]entities.Tag)
	for _, tag := range tags {
		tagsMap[tag.ID] = tag
	}

	product, err := uc.repo.Product().Create(ctx, entities.Product{
		Name:       input.Name,
		Priority:   input.Priority,
		CategoryID: input.CategoryID,
		Status:     entities.StatusProductNeedBuying,
	})
	if err != nil {
		return entities.FullProduct{}, err
	}

	productTags := []entities.Tag{}
	for _, tagID := range input.TagsID {
		tag, ok := tagsMap[tagID]
		if !ok {
			return entities.FullProduct{}, errorsStatus.New(http.StatusBadRequest, "bruh", "нет нек-ых тегов с данными idшниками")
		}

		if err := uc.repo.Grouper().AddTagInProduct(ctx, product.ID, tagID); err != nil {
			return entities.FullProduct{}, err
		}

		productTags = append(productTags, tag)
	}

	category, err := uc.repo.Grouper().CategoryFindByID(ctx, product.CategoryID)
	return entities.FullProduct{
		Product:  product,
		Tags:     productTags,
		Category: category,
	}, err
}
