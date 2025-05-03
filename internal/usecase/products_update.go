package usecase

import (
	"context"
	"glbackend/internal/entities"
	"glbackend/internal/repo"
	"glbackend/internal/utils"
	"slices"
	"time"
)

type (
	ProductsUpdateUsecase interface {
		Execute(ctx context.Context, input ProductsUpdateInput) (entities.FullProduct, error)
	}

	productsUpdateUsecase struct {
		repo       repo.Repo
		ctxTimeout time.Duration
	}

	ProductsUpdateInput struct {
		ID         uint                   `json:"-"`
		Name       string                 `json:"name"`
		Priority   int                    `json:"priority"`
		Icon       string                 `json:"icon"`
		Status     entities.StatusProduct `json:"status"`
		CategoryID uint                   `json:"category_id"`
		TagsID     []uint                 `json:"tags_id"`
	}
)

func NewProductsUpdateUsecase(
	repo repo.Repo,
	timeout time.Duration,
) ProductsUpdateUsecase {
	return &productsUpdateUsecase{
		repo:       repo,
		ctxTimeout: timeout,
	}
}

func (uc *productsUpdateUsecase) Execute(ctx context.Context, input ProductsUpdateInput) (entities.FullProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	product, err := uc.repo.Product().FindByID(ctx, input.ID)
	if err != nil {
		return entities.FullProduct{}, err
	}

	if err = utils.CopyFields(&input, &product); err != nil {
		return entities.FullProduct{}, err
	}

	product, err = uc.repo.Product().Update(ctx, product)
	if err != nil {
		return entities.FullProduct{}, err
	}

	// TODO: add transaction
	productTags, err := uc.repo.Grouper().TagFindAllByProductID(ctx, product.ID)
	if err != nil {
		return entities.FullProduct{}, err
	}

	for _, tagID := range input.TagsID {
		idx := slices.IndexFunc(productTags, func(tag entities.Tag) bool {
			return tag.ID == tagID
		})

		if idx < 0 {
			if err := uc.repo.Grouper().AddTagInProduct(ctx, product.ID, tagID); err != nil {
				return entities.FullProduct{}, err
			}
		}
	}

	for _, tag := range productTags {
		idx := slices.IndexFunc(input.TagsID, func(tagID uint) bool {
			return tag.ID == tagID
		})

		if idx < 0 {
			if err := uc.repo.Grouper().RemoveTagInProduct(ctx, product.ID, tag.ID); err != nil {
				return entities.FullProduct{}, err
			}
		}
	}

	tags, err := uc.repo.Grouper().TagFindAllByProductID(ctx, product.ID)
	if err != nil {
		return entities.FullProduct{}, err
	}
	category, err := uc.repo.Grouper().CategoryFindByID(ctx, product.CategoryID)
	if err != nil {
		return entities.FullProduct{}, err
	}
	return entities.FullProduct{
		Product:  product,
		Tags:     tags,
		Category: category,
	}, nil
}
