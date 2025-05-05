package entities

import (
	"context"
)

type StatusProduct string

const (
	StatusProductBuying        StatusProduct = "bought"
	StatusProductNeedBuying    StatusProduct = "need_buying"
	StatusProductNotNeedBuying StatusProduct = "not_need_buying"
)

type Product struct {
	ID         uint          `json:"id"`
	Name       string        `json:"name"`
	Priority   int           `json:"priority"`
	Status     StatusProduct `json:"status"`
	CategoryID uint          `json:"category_id"`
	Icon       string        `json:"icon"`
}

type FullProduct struct {
	Product  `json:",inline"`
	Tags     []Tag    `json:"tags"`
	Category Category `json:"category"`
}

type ProductRepository interface {
	Create(ctx context.Context, product Product) (Product, error)
	Update(ctx context.Context, product Product) (Product, error)
	Delete(ctx context.Context, id uint) error
	FindAll(ctx context.Context) ([]Product, error)
	FindByID(ctx context.Context, id uint) (Product, error)
	FindAllByTagID(ctx context.Context, tagID uint) ([]Product, error)
	FindAllByIDs(ctx context.Context, ids []uint) ([]Product, error)
}
