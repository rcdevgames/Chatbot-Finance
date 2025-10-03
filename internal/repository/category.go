package repository

import (
	"chatbot/internal/model"
)

type CategoryRepository struct {
	db *SupabaseClient
}

func NewCategoryRepository(db *SupabaseClient) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAllCategories() ([]model.Category, error) {
	var categories []model.Category
	params := map[string]interface{}{
		"order": "type.asc,name.asc",
	}
	err := r.db.Get("categories", params, &categories)
	if err != nil {
		return nil, HandleSupabaseError(err)
	}

	return categories, nil
}

func (r *CategoryRepository) GetCategoriesByType(categoryType string) ([]model.Category, error) {
	var categories []model.Category
	params := map[string]interface{}{
		"type":  categoryType,
		"order": "name.asc",
	}
	err := r.db.Get("categories", params, &categories)
	if err != nil {
		return nil, HandleSupabaseError(err)
	}

	return categories, nil
}

func (r *CategoryRepository) GetCategoryByName(name string) (*model.Category, error) {
	var categories []model.Category
	params := map[string]interface{}{
		"name":  name,
		"limit": 1,
	}
	err := r.db.Get("categories", params, &categories)
	if err != nil {
		return nil, HandleSupabaseError(err)
	}

	if len(categories) == 0 {
		return nil, nil
	}

	return &categories[0], nil
}

func (r *CategoryRepository) SearchCategoryByKeyword(keyword string) ([]model.Category, error) {
	var categories []model.Category
	// Note: Contains might not work with custom client, implement simple filtering
	err := r.db.Get("categories", nil, &categories)
	if err != nil {
		return nil, HandleSupabaseError(err)
	}

	// Filter by keywords (simple implementation)
	var filtered []model.Category
	for _, cat := range categories {
		for _, kw := range cat.Keywords {
			if kw == keyword {
				filtered = append(filtered, cat)
				break
			}
		}
	}

	return filtered, nil
}