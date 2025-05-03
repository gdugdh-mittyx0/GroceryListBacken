package entities

import (
	"context"
)

type Tag struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Category struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type GrouperRepository interface {
	TagCreate(ctx context.Context, tag Tag) (Tag, error)
	TagUpdate(ctx context.Context, tag Tag) (Tag, error)
	TagDelete(ctx context.Context, id uint) error
	TagFindAll(ctx context.Context) ([]Tag, error)
	TagFindByID(ctx context.Context, id uint) (Tag, error)
	TagFindByName(ctx context.Context, name string) (Tag, error)
	TagFindAllByProductID(ctx context.Context, productID uint) ([]Tag, error)

	AddTagInProduct(ctx context.Context, productID, tagID uint) error
	RemoveTagInProduct(ctx context.Context, productID, tagID uint) error

	CategoryCreate(ctx context.Context, category Category) (Category, error)
	CategoryUpdate(ctx context.Context, category Category) (Category, error)
	CategoryDelete(ctx context.Context, id uint) error
	CategoryFindAll(ctx context.Context) ([]Category, error)
	CategoryFindByID(ctx context.Context, id uint) (Category, error)
	CategoryFindByName(ctx context.Context, name string) (Category, error)
}
