package repo

import (
	"context"
	"errors"
	"glbackend/internal/entities"

	"gorm.io/gorm"
)

type ProductGorm struct {
	ID         uint                   `gorm:"column:id"`
	Name       string                 `gorm:"column:name"`
	Priority   int                    `gorm:"column:priority;default:0"`
	Status     entities.StatusProduct `gorm:"column:status"`
	CategoryID uint                   `gorm:"column:category_id"`
	Icon       string                 `gorm:"column:icon"`
}

func (product ProductGorm) TableName() string {
	return "product"
}

func toProductGorm(product entities.Product) ProductGorm {
	return ProductGorm{
		ID:         product.ID,
		Name:       product.Name,
		Priority:   product.Priority,
		Status:     product.Status,
		CategoryID: product.CategoryID,
		Icon:       product.Icon,
	}
}

func (productGorm ProductGorm) toEntity() entities.Product {
	return entities.Product{
		ID:         productGorm.ID,
		Name:       productGorm.Name,
		Priority:   productGorm.Priority,
		Status:     productGorm.Status,
		CategoryID: productGorm.CategoryID,
		Icon:       productGorm.Icon,
	}
}

type ProductGSQL struct {
	db        GSQL
	tableName string
}

func NewProductGSQL(db GSQL) entities.ProductRepository {
	return &ProductGSQL{
		db:        db,
		tableName: "product",
	}
}

func (r ProductGSQL) Create(ctx context.Context, product entities.Product) (entities.Product, error) {
	productGorm := toProductGorm(product)
	if err := r.db.Create(ctx, r.tableName, &productGorm); err != nil {
		return entities.Product{}, err
	}
	return productGorm.toEntity(), nil
}

func (r ProductGSQL) Delete(ctx context.Context, id uint) error {
	if _, err := r.FindByID(ctx, id); err != nil {
		return err
	}
	return r.db.Delete(ctx, r.tableName, &entities.Product{ID: id}, &entities.Product{})
}

func (r ProductGSQL) Update(ctx context.Context, product entities.Product) (entities.Product, error) {
	updatesMap := map[string]interface{}{
		"name":        product.Name,
		"priority":    product.Priority,
		"status":      product.Status,
		"category_id": product.CategoryID,
		"icon":        product.Icon,
	}

	if err := r.db.UpdateOne(ctx, r.tableName, &updatesMap, "id = ?", product.ID); err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return entities.Product{}, errors.New("product_not_found")
		default:
			return entities.Product{}, errors.New("error_update_product")
		}
	}
	return product, nil
}

func (r ProductGSQL) FindAll(ctx context.Context) ([]entities.Product, error) {
	var productsGorm []ProductGorm

	if err := r.db.BeginFind(ctx, r.tableName).Find(&productsGorm); err != nil {
		return nil, errors.New("error_find_products")
	}

	products := make([]entities.Product, 0, len(productsGorm))
	for _, productGorm := range productsGorm {
		products = append(products, productGorm.toEntity())
	}

	return products, nil
}

func (r ProductGSQL) FindByID(ctx context.Context, id uint) (entities.Product, error) {
	var productGorm ProductGorm
	if err := r.db.BeginFind(ctx, r.tableName).Where("id = ?", id).First(&productGorm); err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return entities.Product{}, err
		default:
			return entities.Product{}, errors.New("error_find_product")
		}
	}
	return productGorm.toEntity(), nil
}

func (r ProductGSQL) FindAllByTagID(ctx context.Context, tagID uint) ([]entities.Product, error) {
	var productsGorm []ProductGorm

	find := r.db.BeginFind(ctx, r.tableName).
		Select(`product.*`).
		Join(`LEFT JOIN taginproduct as tip ON tip.product_id = product.id`).
		Where(`tip.tag_id = ?`, tagID)

	if err := find.Find(&productsGorm); err != nil {
		return nil, errors.New("error_find_products")
	}

	products := make([]entities.Product, 0, len(productsGorm))
	for _, productGorm := range productsGorm {
		products = append(products, productGorm.toEntity())
	}

	return products, nil
}
