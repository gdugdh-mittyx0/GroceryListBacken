package repo

import (
	"context"
	"errors"
	"glbackend/internal/entities"

	"gorm.io/gorm"
)

type TagGorm struct {
	ID    uint   `gorm:"column:id"`
	Name  string `gorm:"column:name"`
	Color string `gorm:"column:color"`
}

func (tag TagGorm) TableName() string {
	return "tag"
}

func toTagGorm(tag entities.Tag) TagGorm {
	return TagGorm{
		ID:    tag.ID,
		Name:  tag.Name,
		Color: tag.Color,
	}
}

func (tagGorm TagGorm) toEntity() entities.Tag {
	return entities.Tag{
		ID:    tagGorm.ID,
		Name:  tagGorm.Name,
		Color: tagGorm.Color,
	}
}

type CategoryGorm struct {
	ID    uint   `gorm:"column:id"`
	Name  string `gorm:"column:name"`
	Color string `gorm:"column:color"`
}

func (category CategoryGorm) TableName() string {
	return "category"
}

func toCategoryGorm(category entities.Category) CategoryGorm {
	return CategoryGorm{
		ID:    category.ID,
		Name:  category.Name,
		Color: category.Color,
	}
}

func (categoryGorm CategoryGorm) toEntity() entities.Category {
	return entities.Category{
		ID:    categoryGorm.ID,
		Name:  categoryGorm.Name,
		Color: categoryGorm.Color,
	}
}

type TagInProductGorm struct {
	ID        uint `gorm:"column:id"`
	ProductID uint `gorm:"column:product_id"`
	TagID     uint `gorm:"column:tag_id"`
}

func (tagInProduct TagInProductGorm) TableName() string {
	return "taginproduct"
}

type GrouperGSQL struct {
	db                    GSQL
	tagTableName          string
	categoryTableName     string
	tagInProductTableName string
}

func NewGrouperGSQL(db GSQL) entities.GrouperRepository {
	return &GrouperGSQL{
		db:                    db,
		tagTableName:          "tag",
		categoryTableName:     "category",
		tagInProductTableName: "taginproduct",
	}
}

func (r GrouperGSQL) TagCreate(ctx context.Context, tag entities.Tag) (entities.Tag, error) {
	tagGorm := toTagGorm(tag)
	if err := r.db.Create(ctx, r.tagTableName, &tagGorm); err != nil {
		return entities.Tag{}, err
	}
	return tagGorm.toEntity(), nil
}

func (r GrouperGSQL) TagUpdate(ctx context.Context, tag entities.Tag) (entities.Tag, error) {
	updatesMap := map[string]interface{}{
		"name":  tag.Name,
		"color": tag.Color,
	}

	if err := r.db.UpdateOne(ctx, r.tagTableName, &updatesMap, "id = ?", tag.ID); err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return entities.Tag{}, errors.New("tag_not_found")
		default:
			return entities.Tag{}, errors.New("error_update_tag")
		}
	}
	return tag, nil
}

func (r GrouperGSQL) TagDelete(ctx context.Context, id uint) error {
	return r.db.Delete(ctx, r.tagTableName, &entities.Tag{ID: id}, &entities.Tag{})
}

func (r GrouperGSQL) TagFindAll(ctx context.Context) ([]entities.Tag, error) {
	var tagsGorm []TagGorm
	if err := r.db.BeginFind(ctx, r.tagTableName).Find(&tagsGorm); err != nil {
		return nil, errors.New("error_find_tags")
	}

	tags := make([]entities.Tag, 0, len(tagsGorm))
	for _, tagGorm := range tagsGorm {
		tags = append(tags, tagGorm.toEntity())
	}

	return tags, nil
}

func (r GrouperGSQL) TagFindByID(ctx context.Context, id uint) (entities.Tag, error) {
	var tagGorm TagGorm
	if err := r.db.BeginFind(ctx, r.tagTableName).Where("id = ?", id).First(&tagGorm); err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return entities.Tag{}, err
		default:
			return entities.Tag{}, errors.New("error_find_tag")
		}
	}
	return tagGorm.toEntity(), nil
}

func (r GrouperGSQL) TagFindAllByProductID(ctx context.Context, productID uint) ([]entities.Tag, error) {
	var tagsGorm []TagGorm

	find := r.db.BeginFind(ctx, r.tagTableName).
		Select(`tag.*`).
		Join(`LEFT JOIN taginproduct as tip ON tip.tag_id = tag.id`).
		Where(`tip.product_id = ?`, productID)

	if err := find.Find(&tagsGorm); err != nil {
		return nil, errors.New("error_find_tags")
	}

	tags := make([]entities.Tag, 0, len(tagsGorm))
	for _, tagGorm := range tagsGorm {
		tags = append(tags, tagGorm.toEntity())
	}

	return tags, nil
}

func (r GrouperGSQL) TagFindByName(ctx context.Context, name string) (entities.Tag, error) {
	var tagGorm TagGorm
	if err := r.db.BeginFind(ctx, r.tagTableName).Where("name = ?", name).First(&tagGorm); err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return entities.Tag{}, err
		default:
			return entities.Tag{}, errors.New("error_find_tag")
		}
	}
	return tagGorm.toEntity(), nil
}

func (r GrouperGSQL) AddTagInProduct(ctx context.Context, productID, tagID uint) error {
	if err := r.db.Create(ctx, r.tagInProductTableName, &TagInProductGorm{
		TagID:     tagID,
		ProductID: productID,
	}); err != nil {
		return errors.New("error_add_tag_in_product")
	}
	return nil
}

func (r GrouperGSQL) RemoveTagInProduct(ctx context.Context, productID, tagID uint) error {
	if err := r.db.Delete(ctx, r.tagInProductTableName, nil, "tag_id = ? AND product_id = ?", tagID, productID); err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return errors.New("tag_in_product_not_found")
		default:
			return errors.New("error_remove_tag_from_product")
		}
	}
	return nil
}

func (r GrouperGSQL) CategoryCreate(ctx context.Context, category entities.Category) (entities.Category, error) {
	categoryGorm := toCategoryGorm(category)
	if err := r.db.Create(ctx, r.categoryTableName, &categoryGorm); err != nil {
		return entities.Category{}, err
	}
	return categoryGorm.toEntity(), nil
}

func (r GrouperGSQL) CategoryUpdate(ctx context.Context, category entities.Category) (entities.Category, error) {
	updatesMap := map[string]interface{}{
		"name":  category.Name,
		"color": category.Color,
	}

	if err := r.db.UpdateOne(ctx, r.categoryTableName, &updatesMap, "id = ?", category.ID); err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return entities.Category{}, errors.New("category_not_found")
		default:
			return entities.Category{}, errors.New("error_update_category")
		}
	}
	return category, nil
}

func (r GrouperGSQL) CategoryDelete(ctx context.Context, id uint) error {
	return r.db.Delete(ctx, r.categoryTableName, &entities.Category{ID: id}, &entities.Category{})
}

func (r GrouperGSQL) CategoryFindAll(ctx context.Context) ([]entities.Category, error) {
	var categoriesGorm []CategoryGorm
	if err := r.db.BeginFind(ctx, r.categoryTableName).Find(&categoriesGorm); err != nil {
		return nil, errors.New("error_find_categories")
	}

	categories := make([]entities.Category, 0, len(categoriesGorm))
	for _, categoryGorm := range categoriesGorm {
		categories = append(categories, categoryGorm.toEntity())
	}

	return categories, nil
}

func (r GrouperGSQL) CategoryFindByID(ctx context.Context, id uint) (entities.Category, error) {
	var categoryGorm CategoryGorm
	if err := r.db.BeginFind(ctx, r.categoryTableName).Where("id = ?", id).First(&categoryGorm); err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return entities.Category{}, err
		default:
			return entities.Category{}, errors.New("error_find_category")
		}
	}
	return categoryGorm.toEntity(), nil
}

func (r GrouperGSQL) CategoryFindByName(ctx context.Context, name string) (entities.Category, error) {
	var categoryGorm CategoryGorm
	if err := r.db.BeginFind(ctx, r.categoryTableName).Where("name = ?", name).First(&categoryGorm); err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return entities.Category{}, err
		default:
			return entities.Category{}, errors.New("error_find_category")
		}
	}
	return categoryGorm.toEntity(), nil
}
