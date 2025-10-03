package repository

import (
	"chatbot/internal/model"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAllCategories() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Order("type ASC, name ASC").Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) GetCategoriesByType(categoryType string) ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Where("type = ?", categoryType).
		Order("name ASC").
		Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) GetCategoryByName(name string) (*model.Category, error) {
	var category model.Category
	err := r.db.Where("name = ?", name).First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) SearchCategoryByKeyword(keyword string) ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Where("keywords LIKE ?", "%"+keyword+"%").
		Order("type ASC, name ASC").
		Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) CreateCategory(category *model.Category) error {
	return r.db.Create(category).Error
}

func (r *CategoryRepository) UpdateCategory(category *model.Category) error {
	return r.db.Save(category).Error
}

func (r *CategoryRepository) DeleteCategory(categoryID string) error {
	return r.db.Delete(&model.Category{}, "id = ?", categoryID).Error
}